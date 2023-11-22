package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/environnement"
	"vivarium/organisme"
	"vivarium/terrain"
)

/* ================================================ new server ===================================================== */

var (
	idCount        int = 0 // Global variable, used to add id to each new creature
	ecosystem      *environnement.Environment
	terr           *terrain.Terrain
	clients        = make(map[*websocket.Conn]bool)
	mutex          sync.RWMutex
	ecosystemMutex = &sync.RWMutex{} // Used to protect ecosystem resources
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(terrain *terrain.Terrain, ecosystem *environnement.Environment, w http.ResponseWriter, r *http.Request) {
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

	go func() {
		defer ws.Close()
		for {
			// Read message from client
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}

			// 解析收到的 JSON 消息
			var data map[string]interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}
			fmt.Printf("Parsed JSON: %v\n", data)

			// 根据类型处理生物添加请求
			switch data["type"] {
			case "plant":
				//fmt.Println("植物！")
				handleAddPlantRequest(data, ecosystem, terrain)
			case "insecte":
				//fmt.Println("昆虫！")
				handleAddInsectRequest(data, ecosystem, terrain)
			}
		}
	}()
}

// Handle requests to add plants
func handleAddPlantRequest(data map[string]interface{}, env *environnement.Environment, t *terrain.Terrain) {
	// Extract plant data from the request
	plantTypeStr := data["plantType"].(string)
	// 使用映射进行转换
	plantType, exists := enums.StringToMyEspece[plantTypeStr]
	if !exists {
		log.Printf("Invalid plant type: %s", plantTypeStr)
		return
	}
	posXStr := data["posX"].(string)
	posX, err := strconv.Atoi(posXStr)
	posYStr := data["posY"].(string)
	posY, err := strconv.Atoi(posYStr)
	ageStr := data["plantAge"].(string)
	age, err := strconv.Atoi(ageStr)
	etatSanteStr := data["etatSante"].(string)
	etatSante, err := strconv.Atoi(etatSanteStr)
	if err != nil {
		// handle error
	}

	// Create new plant and add it to the environment and terrain
	// 使用锁保护生态系统状态
	ecosystemMutex.Lock()
	defer ecosystemMutex.Unlock()

	newPlant := organisme.NewPlante(idCount, age, posX, posY, etatSante, plantType)
	idCount++
	env.AjouterOrganisme(newPlant)
	t.AddOrganism(newPlant.GetID(), newPlant.Espece.String(), posX, posY)
}

// Handle requests to add insectes
func handleAddInsectRequest(data map[string]interface{}, env *environnement.Environment, t *terrain.Terrain) {
	// Extract plant data from the request
	insecteTypeStr := data["insecteType"].(string)
	// Convert using mapping
	insecteType, exists := enums.StringToMyEspece[insecteTypeStr]
	if !exists {
		log.Printf("Invalid plant type: %s", insecteTypeStr)
		return
	}
	fmt.Println("昆虫的data：", data)
	posXStr := data["posX"].(string)
	posX, err := strconv.Atoi(posXStr)
	posYStr := data["posY"].(string)
	posY, err := strconv.Atoi(posYStr)
	ageStr := data["insecteAge"].(string)
	age, err := strconv.Atoi(ageStr)
	vitesseStr := data["vitesse"].(string)
	vitesse, err := strconv.Atoi(vitesseStr)
	energyStr := data["energy"].(string)
	energy, err := strconv.Atoi(energyStr)
	// capaciteReproductionStr := data["capaciteReproduction"].(string)
	// capaciteReproduction, err := strconv.Atoi(capaciteReproductionStr)
	// niveauFaimStr := data["niveauFaim"].(string)
	// niveauFaim, err := strconv.Atoi(niveauFaimStr)
	sexeStr := data["sexe"].(string)
	sexe, _ := enums.StringToSexe[sexeStr]
	envieReproduireStr := data["envieReproduire"].(string)
	envieReproduire := false
	if envieReproduireStr == "true" {
		envieReproduire = true
	}
	if err != nil {
		return
	}

	// Create new plant and add it to the environment and terrain
	// 使用锁保护生态系统状态
	ecosystemMutex.Lock()
	defer ecosystemMutex.Unlock()

	newInsecte := organisme.NewInsecte(idCount, age, posX, posY, vitesse, energy, sexe, insecteType, envieReproduire)
	idCount++
	t.AddOrganism(newInsecte.GetID(), newInsecte.Espece.String(), posX, posY)
	env.AjouterOrganisme(newInsecte)
}

func main() {
	// 初始化生态系统
	newEcosystem, newTerrain, newId := environnement.InitializeEcosystem(idCount)
	ecosystem = newEcosystem
	terr = newTerrain
	idCount = newId

	// 启动生态模拟
	go func() {
		ticker := time.NewTicker(time.Second * 1) // 每秒更新一次
		for {
			<-ticker.C

			ecosystemMutex.RLock()
			allOrganismes := ecosystem.GetAllOrganisms()
			ecosystemMutex.RUnlock()

			//fmt.Println("所有生物：", allOrganismes)

			for _, org := range allOrganismes {
				go simulateOrganism(org, allOrganismes)
			}

			if terr == nil {
				fmt.Println("terr is nil")
				return
			}

			// 发送更新后的 Terrain 数据给所有客户端
			updateAndSendTerrain(terr)
		}
	}()

	// Set up WebSocket routing
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(terr, ecosystem, w, r)
	})

	// Start the server
	log.Println("WebSocket server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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

func simulateOrganism(org organisme.Organisme, allOrganismes []organisme.Organisme) {
	//fmt.Println("生物", org.GetID(), org.GetEspece())
	switch o := org.(type) {
	case *organisme.Insecte:
		simulateInsecte(o, allOrganismes)
	case *organisme.Plante:
		simulatePlante(o, allOrganismes, *ecosystem.Climat)
	}
}

func simulateInsecte(ins *organisme.Insecte, allOrganismes []organisme.Organisme) {
	//fmt.Println("昆虫线程启动", ins.GetID())

	// 确认 terr 不是 nil
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	// 判断并执行 Manger
	if ins.AFaim() {
		targetEaten := ins.Manger(allOrganismes, terr)
		if targetEaten != nil {
			ecosystemMutex.Lock()
			ecosystem.RetirerOrganisme(targetEaten)
			ecosystemMutex.Unlock()
		}
	}

	// 判断并执行 SeReproduire，更新昆虫的繁殖意愿
	ins.AvoirEnvieReproduire()
	//fmt.Println("得操了：", ins.GetID(), ins.EnvieReproduire)
	// 执行 SeReproduire 行为
	nbBebes, newOrganismes := ins.SeReproduire(allOrganismes, terr)
	// 如果干出东西了，需要更新到服务器中并显示在terrain上
	if nbBebes != 0 {
		ecosystemMutex.Lock()
		for _, newOrg := range newOrganismes {
			newOrg.SetID(idCount)
			idCount++
			ecosystem.AjouterOrganisme(newOrg)
			terr.AddOrganism(newOrg.GetID(), newOrg.GetEspece().String(), newOrg.GetPosX(), newOrg.GetPosY())
		}
		ecosystemMutex.Unlock()
	}

	// 执行 SeDeplacer
	etatOrganisme := ins.CheckEtat(terr)
	if etatOrganisme != nil {
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(etatOrganisme)
		ecosystemMutex.Unlock()
	} else {
		ecosystemMutex.Lock()
		ins.SeDeplacer(terr) // 需要实现 SeDeplacer 方法
		//updateAndSendTerrain(terr)
		ecosystemMutex.Unlock()

		// 执行 SeBattreRandom
		ins.SeBattreRandom(allOrganismes, terr)
	}

	// 执行 Vieillir
	ins.Vieillir(terr)
	if ins.GetAge() > enums.SpeciesAttributes[ins.GetEspece()].MaxAge {
		// The organism reaches its maximum lifespan and should die
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(ins)
		ecosystemMutex.Unlock()
	}
}

func simulatePlante(pl *organisme.Plante, allOrganismes []organisme.Organisme, climat climat.Climat) {
	//fmt.Println("植物线程启动", pl.GetID())

	// 确认 terr 不是 nil
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	// 更新植物的健康状况
	pl.MisaAJour_EtatSante(climat)
	//fmt.Println("更新完成！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！")

	// 检查植物的当前状态
	etatOrganisme := pl.CheckEtat(terr)
	//fmt.Println("植物状态！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！", etatOrganisme)
	if etatOrganisme != nil {
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(etatOrganisme)
		ecosystemMutex.Unlock()
	}

	// 如果植物还活着，尝试繁殖
	//if etatOrganisme == nil {
	//	nbBebes, newOrganismes := pl.Reproduire(allOrganismes, terr)
	//	if nbBebes != 0 {
	//		ecosystemMutex.Lock()
	//		for _, newOrg := range newOrganismes {
	//			newOrg.SetID(idCount)
	//			idCount++
	//			ecosystem.AjouterOrganisme(newOrg)
	//
	//			terr.AddOrganism(newOrg.GetID(), newOrg.GetEspece().String(), newOrg.GetPosX(), newOrg.GetPosY())
	//		}
	//		ecosystemMutex.Unlock()
	//	}
	//}

	// 模拟植物的生命周期
	pl.Vieillir(terr)
	if pl.GetAge() > enums.SpeciesAttributes[pl.GetEspece()].MaxAge {
		// The organism reaches its maximum lifespan and should die
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(pl)
		ecosystemMutex.Unlock()
	}
}
