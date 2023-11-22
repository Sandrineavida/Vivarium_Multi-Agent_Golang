package test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"vivarium/enums"
	"vivarium/environnement"
	"vivarium/organisme"
	"vivarium/terrain"
)

/* ================================================ old server ====================================================== */

// Global variable, used to add id to each new creature
var idCount int = 0

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool) // 连接集合
var mutex = &sync.Mutex{}                    // 用于保护连接集合

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
	newInsecte := organisme.NewInsecte(idCount, age, posX, posY, vitesse, energy, sexe, insecteType, envieReproduire)
	idCount++
	t.AddOrganism(newInsecte.GetID(), newInsecte.Espece.String(), posX, posY)
	time.Sleep(1)
	env.AjouterOrganisme(newInsecte)
	environnement.Insects = append(environnement.Insects, newInsecte)
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
	ecosystem, terrain, newId := environnement.InitializeEcosystem(idCount)

	idCount = newId

	fmt.Println(ecosystem)
	fmt.Println(terrain)

	go func() {
		ticker := time.NewTicker(time.Second * 1) // updated every second
		for {
			<-ticker.C

			allOrganismes := ecosystem.GetAllOrganisms()
			fmt.Println("All allOrganismes:", allOrganismes)

			for _, organisme := range allOrganismes {
				// Check whether organizationme is of Insecte type
				for _, insect := range environnement.Insects {
					if organisme.GetID() == insect.GetID() {
						if insect.AFaim() {
							targetEaten := insect.Manger(allOrganismes, terrain)
							if targetEaten != nil {
								ecosystem.RetirerOrganisme(targetEaten)
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
				}

				// Update the creature's age
				organisme.Vieillir(terrain)
				if organisme.GetAge() > enums.SpeciesAttributes[organisme.GetEspece()].MaxAge {
					// The organism reaches its maximum lifespan and should die
					ecosystem.RetirerOrganisme(organisme)
				}
			}

			// Send updated Terrain data to all clients
			updateAndSendTerrain(terrain)
		}
	}()

}
