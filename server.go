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

// Global variable, used to add id to each new creature
var idCount int = 0

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool) // 连接集合
var mutex = &sync.Mutex{}                    // 用于保护连接集合

func handleConnections(terrain *terrain.Terrain, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer ws.Close()

	// Add new connection to collection
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
	genre enums.MyEspece
	sexe  enums.Sexe
}

// updateAndSendTerrain sends updated Terrain data to all WebSocket clients
func updateAndSendTerrain(t *terrain.Terrain) {
	terrainJSON, err := json.Marshal(t)
	if err != nil {
		log.Println("Error marshalling terrain:", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, terrainJSON)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func main() {

	// Initialize the ecosystem
	ecosystem, terrain := environnement.InitializeEcosystem(idCount)

	fmt.Println(ecosystem)
	fmt.Println(terrain)

	// Update and send Terrain data regularly
	go func() {
		ticker := time.NewTicker(time.Second * 2) // updated every second
		for {
			<-ticker.C

			allOrganismes := ecosystem.GetAllOrganisms()
			fmt.Println("操你妈", allOrganismes)

			for _, insect := range environnement.Insects {
				if insect.NiveauFaim >= 6 {
					targetInsecte := insect.Manger(allOrganismes, terrain)
					if targetInsecte != nil {
						ecosystem.RetirerOrganisme(targetInsecte)
					}
				}

				// Check the status of insects
				etatOrganisme := insect.CheckEtat(terrain)
				if etatOrganisme != nil {
					ecosystem.RetirerOrganisme(etatOrganisme)
				} else {
					// Update insect location
					insect.SeDeplacer(terrain)
				}
			}

			// Send updated Terrain data to all clients
			updateAndSendTerrain(terrain)
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
