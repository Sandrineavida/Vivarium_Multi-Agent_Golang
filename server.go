package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
	"vivarium/enums"
	"vivarium/environnement"
	"vivarium/terrain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool) // 连接集合
var mutex = &sync.Mutex{}                    // 用于保护连接集合

func handleConnections(terrain *terrain.Terrain, w http.ResponseWriter, r *http.Request) {
	// 升级初始GET请求到一个websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer ws.Close()

	// 添加新连接到集合
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	// Important: Here, you must send the status of the ecosystem to the client
	// Send biometric information
	infoJSON, err := json.Marshal(terrain)
	if err != nil {
		log.Println("Error marshalling organisms info:", err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, infoJSON)
}

func describeSex(sex enums.Sexe) {
	if sex == enums.Male {
		fmt.Println("The sex is Male.")
	} else if sex == enums.Femelle {
		fmt.Println("The sex is Female.")
	} else {
		fmt.Println("The sex is Hermaphrodite.")
	}
}

type Creature struct {
	genre enums.MyInsect
	sexe  enums.Sexe
}

func main() {
	//sex := enums.Male
	//fmt.Println(sex.String()) // print "Male"
	//
	//cre1 := Creature{
	//	genre: enums.Lombric,
	//	sexe:  enums.Femelle,
	//}
	//describeSex(cre1.sexe)
	//
	//cre2 := Creature{
	//	genre: enums.Escargot,
	//	sexe:  enums.Hermaphrodite,
	//}
	//describeSex(cre2.sexe)

	// Initialize the ecosystem
	ecosystem, terrain := environnement.InitializeEcosystem()

	fmt.Println(ecosystem) // 有问题
	fmt.Println(terrain)

	// 定时更新和发送 Terrain 数据
	go func() {
		ticker := time.NewTicker(time.Second * 1) // updated every second
		for {
			<-ticker.C
			// update insect location
			for _, insect := range environnement.Insects {
				insect.SeDeplacer(terrain)
			}

			// 发送更新后的 Terrain 数据到所有客户端
			terrainJSON, err := json.Marshal(terrain)
			if err != nil {
				log.Println("Error marshalling terrain:", err)
				continue
			}

			mutex.Lock()
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, terrainJSON)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
			mutex.Unlock()
		}
	}()

	// Set up WebSocket routing
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(terrain, w, r)
	})

	// Start the server
	log.Println("WebSocket server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
