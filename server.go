package server

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
	idCount            int = 0 // Global variable, used to add id to each new creature
	ecosystem          *environnement.Environment
	terr               *terrain.Terrain
	clients            = make(map[*websocket.Conn]bool)
	mutex              sync.RWMutex
	ecosystemMutex     = &sync.RWMutex{} // Used to protect ecosystem resources
	wg                 sync.WaitGroup
	EcosystemForEbiten *environnement.Environment
	ecoMutex           sync.RWMutex
	//EcosystemChannel   chan *environnement.Environment
	PauseSignal        = make(chan bool) // 用于暂停和继续的信号
	pauseMutex         sync.RWMutex      // 锁，用于保护 pauseSignal
	isSimulationPaused bool              // 表示仿真是否已经暂停
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
	// etatSanteStr := data["etatSante"].(string)
	// etatSante, err := strconv.Atoi(etatSanteStr)
	if err != nil {
		// handle error
	}

	// Create new plant and add it to the environment and terrain
	// 使用锁保护生态系统状态
	ecosystemMutex.Lock()
	defer ecosystemMutex.Unlock()

	newPlant := organisme.NewPlante(idCount, age, posX, posY, plantType)
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
	//fmt.Println("昆虫的data：", data)
	posXStr := data["posX"].(string)
	posX, err := strconv.Atoi(posXStr)
	posYStr := data["posY"].(string)
	posY, err := strconv.Atoi(posYStr)
	ageStr := data["insecteAge"].(string)
	age, err := strconv.Atoi(ageStr)
	//vitesseStr := data["vitesse"].(string)
	//vitesse, err := strconv.Atoi(vitesseStr)
	// energyStr := data["energy"].(string)
	// energy, err := strconv.Atoi(energyStr)
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

	newInsecte := organisme.NewInsecte(idCount, age, posX, posY, sexe, insecteType, envieReproduire)
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

	// 设置定时器在5秒后将天气改回 Rien
	time.AfterFunc(5*time.Second, func() {
		ecosystemMutex.Lock()
		ecosystem.Climat.ChangerConditions(enums.Rien)
		t.Meteo = enums.Rien
		ecosystemMutex.Unlock()
	})
}

// 控制仿真的暂停和继续的goroutine
func controlSimulation() {
	for {
		select {
		case <-PauseSignal:
			// 切换仿真状态
			isSimulationPaused = !isSimulationPaused
			if isSimulationPaused {
				fmt.Println("Simulation paused")
			} else {
				fmt.Println("Simulation resumed")
			}
		}
	}
}

func StartServer() {
	// 只在最开始设置一次随机数种子 - 2023.11.22
	rand.Seed(time.Now().UnixNano())

	// 该goroutine用于控制仿真的暂停和继续：
	go controlSimulation()

	// 更新 Hour 并根据当前 Hour 更新气候

	// 初始化生态系统
	newEcosystem, newTerrain, newId := environnement.InitializeEcosystem(idCount)
	//EcosystemChannel = make(chan *environnement.Environment)
	EcosystemForEbiten = &environnement.Environment{}
	ecosystem = newEcosystem
	terr = newTerrain
	idCount = newId
	isinit := true

	// 启动生态模拟
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			<-ticker.C

			// 更新 Hour 并根据当前 Hour 更新气候
			ecosystem.Hour = (ecosystem.Hour + 1) % 24
			//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!当前天气：", ecosystem.Climat.Meteo, "当前时间：", ecosystem.Hour)
			ecosystem.Climat.UpdateClimat_24H(ecosystem.Hour, isinit)
			isinit = false

			ecosystemMutex.RLock()
			allOrganismes := ecosystem.GetAllOrganisms()
			// 创建一个新的切片用于存储未死亡的生物
			aliveOrganismes := make([]organisme.Organisme, 0)
			// 创建一个新的切片用于存储未死亡的生物
			for _, org := range allOrganismes {
				if !org.GetEtat() {
					// 如果生物未死亡，添加到新的切片中
					aliveOrganismes = append(aliveOrganismes, org)
				}
			}
			allOrganismes = aliveOrganismes
			ecosystemMutex.RUnlock()

			wg.Add(len(allOrganismes))
			for _, org := range allOrganismes {
				go func(o organisme.Organisme) {
					simulateOrganism(o, allOrganismes)
					wg.Done()
				}(org)
			}

			//fmt.Println("操你妈大草消掉了", allOrganismes, "我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼我是傻逼")

			wg.Wait() // 等待所有 simulateOrganism goroutines 完成
			//updateAndSendTerrain(terr)
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

func simulateOrganism(org organisme.Organisme, allOrganismes []organisme.Organisme) {
	//fmt.Println("生物", org.GetID(), org.GetEspece())
	for {
		// 检查仿真是否已暂停，如果是，则等待
		if isSimulationPaused {
			time.Sleep(time.Millisecond * 100) // 等待一小段时间再检查状态
			continue
		}

		switch o := org.(type) {
		case *organisme.Insecte:
			simulateInsecte(o, allOrganismes, *ecosystem.Climat)
			time.Sleep(time.Millisecond * 100)
		case *organisme.Plante:
			simulatePlante(o, allOrganismes, *ecosystem.Climat)
			time.Sleep(time.Millisecond * 100)
		}

		// EcosystemForEbiten用于提供给ebiten，让所有精灵获取到EcosystemForEbiten中的所有生物信息，并更新精灵
		// 每次每个生物仿真结束后都会对EcosystemForEbiten进行更新
		ecoMutex.Lock() // 开始写操作前加锁
		EcosystemForEbiten = ecosystem
		ecoMutex.Unlock()
		break

		// 发送数据
		// <- EcosystemForEbiten
	}
}

func simulateInsecte(ins *organisme.Insecte, allOrganismes []organisme.Organisme, climat climat.Climat) {
	//fmt.Println("昆虫线程启动", ins.GetID())

	// 确认 terr 不是 nil
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	//fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫开始行动！！！！！:::::::", ins.Energie)

	// hotfix-1124: 先感受一下是不是火灾了 (这样其实新生儿就也可以马上受到火灾影响)
	//if ins.PerceptClimat(climat) {
	//	ins.UpdateEnergie()
	//	// 看看有没有被烧死
	//	burnt_to_death := ins.CheckEtat(terr)
	//	if burnt_to_death != nil {
	//		//ecosystemMutex.Lock()
	//		//ecosystem.RetirerOrganisme(burnt_to_death)
	//		//ecosystemMutex.Unlock()
	//		// fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫【【【被烧死】】】死了！！！！！!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	//		return
	//	}
	//}

	// 感知环境，根据环境的极端程度进行UpdateEnergie
	severity := ins.PerceptClimat(climat)
	ins.UpdateEnergie(severity)
	// 判断并执行 Manger
	if ins.AFaim() {
		//fmt.Println("[", ins.OrganismeID, ins.Espece, "]:  昆虫饿了！！！！！:::::::", ins.Energie)
		ins.Manger(allOrganismes, terr)
		//targetEaten := ins.Manger(allOrganismes, terr)
		//if targetEaten != nil {
		//	//ecosystemMutex.Lock()
		//	//ecosystem.RetirerOrganisme(targetEaten)
		//	//ecosystemMutex.Unlock()
		//}
	}

	severity = ins.PerceptClimat(climat)
	ins.UpdateEnergie(severity)

	// 判断并执行 SeReproduire，更新昆虫的繁殖意愿
	ins.AvoirEnvieReproduire()
	//fmt.Println("我是 ", ins.Espece, "", ins.GetID(), "", ins.EnvieReproduire)
	// 执行 SeReproduire 行为 若果打赢了，就想再找昆虫bang，findTargetAganin为true，但也就找这么一次
	nbBebes, newOrganismes, findTargetAganin := ins.SeReproduire(allOrganismes, terr)
	if findTargetAganin {
		nbBebes, newOrganismes, _ = ins.SeReproduire(allOrganismes, terr)
	}

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

	severity = ins.PerceptClimat(climat)
	ins.UpdateEnergie(severity)

	// 执行 SeDeplacer
	time.Sleep(time.Millisecond * 1000) // 控制每次移动之间的时间间隔
	for i := 0; i < ins.Vitesse; i++ {
		etatOrganisme := ins.CheckEtat(terr)
		if etatOrganisme != nil {
			//ecosystemMutex.Lock()
			//ecosystem.RetirerOrganisme(etatOrganisme)
			//ecosystemMutex.Unlock()
			//fmt.Println("[", ins.OrganismeID, "]: 昆虫死了！")
			return
		} else {
			ecosystemMutex.Lock()
			ins.SeDeplacer(terr) // 执行移动动作
			//updateAndSendTerrain(terr) // 立即更新并发送terrain状态
			ecosystemMutex.Unlock()
		}
		time.Sleep(time.Millisecond * 100) // 控制每次移动之间的时间间隔
	}
	time.Sleep(time.Millisecond * 1000) // 控制每次移动之间的时间间隔

	// 执行 Vieillir
	ins.Vieillir(terr)
	//if ins.GetAge() > enums.SpeciesAttributes[ins.GetEspece()].MaxAge {
	//	// The organism reaches its maximum lifespan and should die
	//	//ecosystemMutex.Lock()
	//	//ecosystem.RetirerOrganisme(ins)
	//	//ecosystemMutex.Unlock()
	//}

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
		//ecosystemMutex.Lock()
		//ecosystem.RetirerOrganisme(etatOrganisme)
		//ecosystemMutex.Unlock()
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
	//if pl.GetAge() > enums.SpeciesAttributes[pl.GetEspece()].MaxAge {
	//	// The organism reaches its maximum lifespan and should die
	//	ecosystemMutex.Lock()
	//	ecosystem.RetirerOrganisme(pl)
	//	ecosystemMutex.Unlock()
	//}

	//updateAndSendTerrain(terr)
}
