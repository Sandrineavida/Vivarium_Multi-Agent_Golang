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

var (
	idCount int = 0 // Global variable, used to add id to each new creature

	ecosystem      *environnement.Environment
	ecosystemMutex sync.RWMutex // Used to protect ecosystem and terrain resources

	EcosystemForEbiten *environnement.Environment
	ecoMutex           sync.RWMutex // Used to protect EcosystemForEbiten resources

	terr    *terrain.Terrain
	clients = make(map[*websocket.Conn]bool)

	wg                   sync.WaitGroup
	mutex                sync.RWMutex // Connection handling for new clients, although in this project only one user will connect via websocket
	aliveOrganismesMutex sync.RWMutex // Used to protect globalAliveOrganismes resources

	PauseSignal        = make(chan bool) // Signal for pause and continue buttons in ebiten
	isSimulationPaused bool              // Indicates whether the simulation has been paused

	globalAliveOrganismes []organisme.Organisme // Drawing for Terrain in html
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(terrain *terrain.Terrain, ecosystem *environnement.Environment, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add new connection to collection
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

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

			// Parse the received JSON message
			var data map[string]interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			switch data["type"] {
			case "plant":
				handleAddPlantRequest(data, ecosystem, terrain)
			case "insecte":
				handleAddInsectRequest(data, ecosystem, terrain)
			case "requestTerrainData":
				updateAndSendTerrain(terrain)
			case "changeMeteo":
				updateMeteoAndSendTerrain(data, terrain)
			}
		}
	}()
}

// Handle requests to add plants
func handleAddPlantRequest(data map[string]interface{}, env *environnement.Environment, t *terrain.Terrain) {
	// Extract plant data from the request
	plantTypeStr := data["plantType"].(string)
	// Convert using mapping
	plantType, exists := enums.StringToMyEspece[plantTypeStr]
	if !exists {
		log.Printf("Invalid plant type: %s", plantTypeStr)
		return
	}
	posXStr := data["posX"].(string)
	posX, _ := strconv.Atoi(posXStr)
	posYStr := data["posY"].(string)
	posY, _ := strconv.Atoi(posYStr)
	ageStr := data["plantAge"].(string)
	age, _ := strconv.Atoi(ageStr)

	// Create new plant and add it to the environment and terrain
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

	posXStr := data["posX"].(string)
	posX, err := strconv.Atoi(posXStr)
	posYStr := data["posY"].(string)
	posY, err := strconv.Atoi(posYStr)
	ageStr := data["insecteAge"].(string)
	age, err := strconv.Atoi(ageStr)
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

	// Clear the current Grid
	for y := range t.Grid {
		for x := range t.Grid[y] {
			t.Grid[y][x] = nil
		}
	}

	// Get the information of living organismes in the current ecological box, and then render it in html
	aliveOrganismesMutex.RLock()
	defer aliveOrganismesMutex.RUnlock()

	// Repopulate Grid based on aliveOrganismes
	for _, org := range globalAliveOrganismes {
		x, y := org.GetPosX(), org.GetPosY()
		t.Grid[y][x] = append(t.Grid[y][x], terrain.CellInfo{
			OrganismID:   org.GetID(),
			OrganismType: org.GetEspece().String(),
		})
	}

	// Update current time
	t.CurrentHour = ecosystem.Hour

	// Update current time Climat and Meteo
	t.Meteo = ecosystem.Climat.Meteo
	t.Luminaire = ecosystem.Climat.Luminaire
	t.Temperature = ecosystem.Climat.Temperature
	t.Humidite = ecosystem.Climat.Humidite
	t.Co2 = ecosystem.Climat.Co2
	t.O2 = ecosystem.Climat.O2

	// Serialize and send updated Terrain information to the client
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

// HTML display for Meteo changes
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

	// Set timer to change weather back to Rien after 5 seconds
	time.AfterFunc(5*time.Second, func() {
		ecosystemMutex.Lock()
		ecosystem.Climat.ChangerConditions(enums.Rien)
		t.Meteo = enums.Rien
		ecosystemMutex.Unlock()
	})
}

// Fonction to control the pause and continuation of simulation
func controlSimulation() {
	for {
		select {
		case <-PauseSignal:
			// Switch simulation state
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
	rand.Seed(time.Now().UnixNano())

	// This goroutine is used to control the pause and continuation of the simulation:
	go controlSimulation()

	// Initialize the ecosystem
	newEcosystem, newTerrain, newId := environnement.InitializeEcosystem(idCount)
	EcosystemForEbiten = &environnement.Environment{}
	ecosystem = newEcosystem
	terr = newTerrain
	idCount = newId
	isinit := true

	// Start ecological simulation
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			<-ticker.C

			// Update Hour and update climate based on current Hour
			ecosystem.Hour = (ecosystem.Hour + 1) % 24
			ecosystem.Climat.UpdateClimat_24H(ecosystem.Hour, isinit)
			isinit = false

			ecosystemMutex.RLock()
			allOrganismes := ecosystem.GetAllOrganisms()
			// Create a new slice to store non-dead creatures
			aliveOrganismes := make([]organisme.Organisme, 0)
			for _, org := range allOrganismes {
				if !org.GetEtat() {
					aliveOrganismes = append(aliveOrganismes, org)
				}
			}
			allOrganismes = aliveOrganismes
			ecosystemMutex.RUnlock()

			aliveOrganismesMutex.Lock()
			globalAliveOrganismes = aliveOrganismes
			aliveOrganismesMutex.Unlock()

			wg.Add(len(allOrganismes))
			for _, org := range allOrganismes {
				go func(o organisme.Organisme) {
					simulateOrganism(o, allOrganismes)
					wg.Done()
				}(org)
			}

			wg.Wait() // Wait for all simulateOrganism goroutines to complete
		}
	}()

	// Set up WebSocket routing
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(terr, ecosystem, w, r)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// Start the server
	log.Println("WebSocket server started on :8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func simulateOrganism(org organisme.Organisme, allOrganismes []organisme.Organisme) {
	for {
		// Check if the simulation is paused, if yes, wait
		if isSimulationPaused {
			time.Sleep(time.Millisecond * 100)
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

		// EcosystemForEbiten is used to be provided to ebiten, so that all spritess can obtain all biological information in EcosystemForEbiten and update the sprites
		// EcosystemForEbiten will be updated every time after each biological simulation is completed
		ecoMutex.Lock()
		EcosystemForEbiten = ecosystem
		ecoMutex.Unlock()
		break
	}
}

// Insect simulation logic
func simulateInsecte(ins *organisme.Insecte, allOrganismes []organisme.Organisme, climat climat.Climat) {
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	// Perceive the environment and perform UpdateEnergie according to the extreme degree of the environment
	severity := ins.PerceptClimat(climat)
	ins.UpdateEnergie(severity)
	// Determine and execute Manger
	if ins.AFaim() {
		ins.Manger(allOrganismes, terr)
	}

	severity = ins.PerceptClimat(climat)
	ins.UpdateEnergie(severity)

	// Determine and execute SeReproduire to update the insect's willingness to reproduce
	ins.AvoirEnvieReproduire()
	// Execute the SeReproduire behavior. If the insect wins the battle with the same sex, it will find the opposite sex here to mate
	nbBebes, newOrganismes, findTargetAganin := ins.SeReproduire(allOrganismes, terr)
	if findTargetAganin {
		nbBebes, newOrganismes, _ = ins.SeReproduire(allOrganismes, terr)
	}

	// Add new insects
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

	// Execute SeDeplacer
	time.Sleep(time.Millisecond * 1000)
	// Depending on the speed of the insect, SeDeplacer is called a number of times matching the speed at the same time
	for i := 0; i < ins.Vitesse; i++ {
		etatOrganisme := ins.CheckEtat(terr)
		if etatOrganisme != nil {
			return
		} else {
			ecosystemMutex.Lock()
			ins.SeDeplacer(terr)
			ecosystemMutex.Unlock()
		}
		time.Sleep(time.Millisecond * 100) // Control the time between each move
	}

	time.Sleep(time.Millisecond * 1000)

	// Execute Vieillir
	ins.Vieillir(terr)
}

// Plant simulation logic
func simulatePlante(pl *organisme.Plante, allOrganismes []organisme.Organisme, climat climat.Climat) {
	if terr == nil {
		fmt.Println("错误: terr")
		return
	}

	// Update plant health
	pl.MisaAJour_EtatSante(climat)
	etatOrganisme := pl.CheckEtat(terr)

	// If the plant is alive, try to propagate it
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

	pl.Vieillir(terr)
}
