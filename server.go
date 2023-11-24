package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/environnement"
	"vivarium/organisme"
	"vivarium/terrain"

	"github.com/gorilla/websocket"
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
			//fmt.Printf("Parsed JSON: %v\n", data)

			// 根据类型处理生物添加请求
			switch data["type"] {
			case "plant":
				//fmt.Println("植物！")
				handleAddPlantRequest(data, ecosystem, terrain)
			case "insecte":
				//fmt.Println("昆虫！")
				handleAddInsectRequest(data, ecosystem, terrain)
			case "requestTerrainData":
				updateAndSendTerrain(terrain)
			case "changeMeteo":
				updateMeteoAndSendTerrain(data, terrain)
				//case "changeClimat":
				//	updateClimatAndSendTerrain(data, terrain)
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
	//vitesseStr := data["vitesse"].(string)
	//vitesse, err := strconv.Atoi(vitesseStr)
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

	newInsecte := organisme.NewInsecte(idCount, age, posX, posY, energy, sexe, insecteType, envieReproduire)
	idCount++
	t.AddOrganism(newInsecte.GetID(), newInsecte.Espece.String(), posX, posY)
	env.AjouterOrganisme(newInsecte)
}

func updateAndSendTerrain(t *terrain.Terrain) {
	mutex.Lock()
	defer mutex.Unlock()

	// 在发送之前更新当前时间
	t.CurrentHour = ecosystem.Hour

	// 更新climat和meteo
	t.Meteo = ecosystem.Climat.Meteo
	t.Luminaire = ecosystem.Climat.Luminaire
	t.Temperature = ecosystem.Climat.Temperature
	t.Humidite = ecosystem.Climat.Humidite
	t.Co2 = ecosystem.Climat.Co2
	t.O2 = ecosystem.Climat.O2

	terrainJSON, err := json.Marshal(t)
	if err != nil {
		log.Println("Error marshalling terrain:", err)
		return
	}

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, terrainJSON)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func updateMeteoAndSendTerrain(data map[string]interface{}, t *terrain.Terrain) {
	meteoTypeStr, ok := data["meteoType"].(string)
	if !ok {
		log.Println("Invalid meteo type")
		return
	}
	meteoType, exists := enums.StringToMeteo[meteoTypeStr]
	if !exists {
		log.Printf("Invalid meteo type: %s", meteoTypeStr)
		return
	}
	ecosystemMutex.Lock()
	ecosystem.Climat.ChangerConditions(meteoType)
	t.Meteo = meteoType
	ecosystemMutex.Unlock()
	//updateAndSendTerrain(terrain)
}

//func updateClimatAndSendTerrain(data map[string]interface{}, t *terrain.Terrain) {
//	luminaireStr, ok := data["Luminaire"].(string)
//	if !ok {
//		log.Println("Invalid luminaire data")
//		return
//	}
//	luminaire, err := strconv.Atoi(luminaireStr)
//	if err != nil {
//		log.Println("Error converting luminaire:", err)
//		return
//	}
//
//	temperatureStr, ok := data["Temperature"].(string)
//	if !ok {
//		log.Println("Invalid temperature data")
//		return
//	}
//	temperature, err := strconv.Atoi(temperatureStr)
//	if err != nil {
//		log.Println("Error converting temperature:", err)
//		return
//	}
//
//	humidite, ok := data["Humidite"].(float32) // JSON中的浮点数默认解析为float64
//	if !ok {
//		log.Println("Invalid humidite data")
//		return
//	}
//
//	// 类似地处理 Co2 和 O2
//	co2, ok := data["Co2"].(float32)
//	if !ok {
//		log.Println("Invalid Co2 data")
//		return
//	}
//
//	o2, ok := data["O2"].(float32)
//	if !ok {
//		log.Println("Invalid O2 data")
//		return
//	}
//
//	ecosystemMutex.Lock()
//	ecosystem.Climat.Luminaire = luminaire
//	ecosystem.Climat.Humidite = humidite
//	ecosystem.Climat.Temperature = temperature
//	ecosystem.Climat.Co2 = co2
//	ecosystem.Climat.O2 = o2
//	ecosystemMutex.Unlock()
//}

var wg sync.WaitGroup

func main() {
	// 只在最开始设置一次随机数种子 - 2023.11.22
	rand.Seed(time.Now().UnixNano())

	// 更新 Hour 并根据当前 Hour 更新气候

	// 初始化生态系统
	newEcosystem, newTerrain, newId := environnement.InitializeEcosystem(idCount)
	ecosystem = newEcosystem
	terr = newTerrain
	idCount = newId
	isinit := true

	// 启动生态模拟
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C

			// 更新 Hour 并根据当前 Hour 更新气候
			ecosystem.Hour = (ecosystem.Hour + 1) % 24
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!当前天气：", ecosystem.Climat.Meteo, "当前时间：", ecosystem.Hour)
			ecosystem.Climat.UpdateClimat_24H(ecosystem.Hour, isinit)
			isinit = false

			ecosystemMutex.RLock()
			allOrganismes := ecosystem.GetAllOrganisms()
			ecosystemMutex.RUnlock()

			wg.Add(len(allOrganismes))
			for _, org := range allOrganismes {
				go func(o organisme.Organisme) {
					simulateOrganism(o, allOrganismes)
					wg.Done()
				}(org)
			}

			wg.Wait() // 等待所有 simulateOrganism goroutines 完成
			//updateAndSendTerrain(terr)
		}
	}()

	// 方法2：管道
	//go func() {
	//	ticker := time.NewTicker(time.Second)
	//	for {
	//		<-ticker.C
	//
	//		doneChan := make(chan struct{})
	//		ecosystemMutex.RLock()
	//		allOrganismes := ecosystem.GetAllOrganisms()
	//		ecosystemMutex.RUnlock()
	//
	//		for _, org := range allOrganismes {
	//			go func(o organisme.Organisme) {
	//				simulateOrganism(o, allOrganismes)
	//				doneChan <- struct{}{}
	//			}(org)
	//		}
	//
	//		for i := 0; i < len(allOrganismes); i++ {
	//			<-doneChan
	//		}
	//		updateAndSendTerrain(terr)
	//	}
	//}()

	// 方法3：Timer定时更新生态系统
	//for {
	//	// 设置一个1秒钟的定时器
	//	timer := time.NewTimer(1 * time.Second)
	//
	//	ecosystemMutex.RLock()
	//	allOrganismes := ecosystem.GetAllOrganisms()
	//	ecosystemMutex.RUnlock()
	//
	//	var wg sync.WaitGroup
	//	for _, org := range allOrganismes {
	//		wg.Add(1)
	//		go func(o organisme.Organisme) {
	//			defer wg.Done()
	//			simulateOrganism(o, allOrganismes)
	//		}(org)
	//	}
	//
	//	// 等待所有goroutines完成
	//	wg.Wait()
	//
	//	if terr != nil {
	//		updateAndSendTerrain(terr)
	//	}
	//
	//	// 等待定时器到时
	//	<-timer.C
	//}

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

func simulateOrganism(org organisme.Organisme, allOrganismes []organisme.Organisme) {
	//fmt.Println("生物", org.GetID(), org.GetEspece())
	switch o := org.(type) {
	case *organisme.Insecte:
		simulateInsecte(o, allOrganismes, *ecosystem.Climat)
		time.Sleep(time.Millisecond * 100)
	case *organisme.Plante:
		simulatePlante(o, allOrganismes, *ecosystem.Climat)
		time.Sleep(time.Millisecond * 100)
	}
}

func simulateInsecte(ins *organisme.Insecte, allOrganismes []organisme.Organisme, climat climat.Climat) {
	//fmt.Println("昆虫线程启动", ins.GetID())

	// 确认 terr 不是 nil
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫开始行动！！！！！:::::::", ins.Energie)

	// hotfix-1124: 先感受一下是不是火灾了 (这样其实新生儿就也可以马上受到火灾影响)
	if ins.PerceptIncendie(climat) {
		ins.UpdateEnergie_Incendie()
		// 看看有没有被烧死
		burnt_to_death := ins.CheckEtat(terr)
		if burnt_to_death != nil {
			ecosystemMutex.Lock()
			ecosystem.RetirerOrganisme(burnt_to_death)
			ecosystemMutex.Unlock()
			fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫【【【被烧死】】】死了！！！！！!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
	}

	// 判断并执行 Manger
	if ins.AFaim() {
		fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫饿了！！！！！:::::::", ins.Energie)
		targetEaten := ins.Manger(allOrganismes, terr)
		if targetEaten != nil {
			ecosystemMutex.Lock()
			ecosystem.RetirerOrganisme(targetEaten)
			ecosystemMutex.Unlock()
		}
	}
	// hotfix-1124: 先感受一下是不是火灾了
	if ins.PerceptIncendie(climat) {
		ins.UpdateEnergie_Incendie()
	}
	etatOrganisme_dead := ins.CheckEtat(terr)
	if etatOrganisme_dead != nil {
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(etatOrganisme_dead)
		ecosystemMutex.Unlock()
		if ins.PerceptIncendie(climat) {
			fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫【【【被烧死】】】死了！！！！！!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		} else {
			fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫@@@@@饿@@@@@死了！！！！！!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		}
		return
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
	// hotfix-1124: 干完了生完了，来烧烧看
	if ins.PerceptIncendie(climat) {
		ins.UpdateEnergie_Incendie()
		// 看看有没有被烧死
		burnt_to_death := ins.CheckEtat(terr)
		if burnt_to_death != nil {
			ecosystemMutex.Lock()
			ecosystem.RetirerOrganisme(burnt_to_death)
			ecosystemMutex.Unlock()
			fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫【【【被烧死】】】死了！！！！！!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
	}

	// 执行 SeDeplacer
	for i := 0; i < ins.Vitesse; i++ {
		etatOrganisme := ins.CheckEtat(terr)
		if etatOrganisme != nil {
			ecosystemMutex.Lock()
			ecosystem.RetirerOrganisme(etatOrganisme)
			ecosystemMutex.Unlock()
			fmt.Println("[", ins.OrganismeID, "]: 昆虫死了！")
			return
		} else {
			ecosystemMutex.Lock()
			ins.SeDeplacer(terr) // 执行移动动作
			//updateAndSendTerrain(terr) // 立即更新并发送terrain状态
			ecosystemMutex.Unlock()
		}
		time.Sleep(time.Millisecond * 100) // 控制每次移动之间的时间间隔
	}

	// 执行 Vieillir
	ins.Vieillir(terr)
	if ins.GetAge() > enums.SpeciesAttributes[ins.GetEspece()].MaxAge {
		// The organism reaches its maximum lifespan and should die
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(ins)
		ecosystemMutex.Unlock()
	}

	//updateAndSendTerrain(terr)
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
	if etatOrganisme != nil {
		fmt.Println("植物要死！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！", etatOrganisme)
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(etatOrganisme)
		ecosystemMutex.Unlock()
	}

	// 如果植物还活着，尝试繁殖
	if etatOrganisme == nil {
		nbBebes, newOrganismes := pl.Reproduire(allOrganismes, terr)
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
	}

	// 模拟植物的生命周期
	pl.Vieillir(terr)
	if pl.GetAge() > enums.SpeciesAttributes[pl.GetEspece()].MaxAge {
		// The organism reaches its maximum lifespan and should die
		ecosystemMutex.Lock()
		ecosystem.RetirerOrganisme(pl)
		ecosystemMutex.Unlock()
	}

	//updateAndSendTerrain(terr)
}
