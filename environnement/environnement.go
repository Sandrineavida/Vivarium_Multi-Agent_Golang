package environnement

import (
	"math/rand"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/organisme"
	"vivarium/terrain"
)

// Environment represents the simulation environment.
type Environment struct {
	Climat     *climat.Climat
	QualiteSol int
	Width      int
	Height     int
	NbPierre   int
	Organismes []organisme.Organisme
}

// NewEnvironment creates a new instance of Environment with default values.
func NewEnvironment(width, height int) *Environment {
	return &Environment{
		Climat:     climat.NewClimat(),
		Width:      width,
		Height:     height,
		Organismes: make([]organisme.Organisme, 0),
		// Set other attributes...
	}
}

// Simuler simulates the environment for a time step.
func (e *Environment) Simuler() {
	// Implementation of simulation step
}

// AjouterOrganisme adds a new organism to the environment.
func (e *Environment) AjouterOrganisme(o organisme.Organisme) {
	// Implementation to add a new organism
	e.Organismes = append(e.Organismes, o)
}

// RetirerOrganisme removes an organism from the environment.
func (e *Environment) RetirerOrganisme(o *organisme.Organisme) {
	// Implementation to remove an organism
	// This might involve searching for the organism in the list and removing it.
}

// MiseAJour updates the environment state.
func (e *Environment) MiseAJour() {
	// Implementation to update the environment
	// This might involve updating the state of each organism, climate changes, etc.
}

/* Written by Zhenyang here */

// Initial number of assumptions
const (
	initialPlantCount  = 10
	initialInsectCount = 5
)

var Insects []*organisme.Insecte

type OrganismeInfo struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Species    string `json:"species"`
	Position_X int    `json:"position_x"`
	Position_Y int    `json:"position_y"`
}

// InitializeEcosystem creates and initializes the environment and creatures of the ecosystem
func InitializeEcosystem() (*Environment, *terrain.Terrain) {
	// Create environment instance
	env := NewEnvironment(10, 10)
	terr := terrain.NewTerrain(10, 10)

	// Add initial plants
	for i := 0; i < initialPlantCount; i++ {
		posX := rand.Intn(10)
		posY := rand.Intn(10)
		plant := organisme.NewPlante(
			i,                                        // ID
			0,                                        // Age
			posX,                                     // positionX
			posY,                                     // positionY
			5,                                        // Rayon
			1,                                        // VitesseDeCroissance
			100,                                      // EtatSante
			1,                                        // Adaptabilite
			enums.ModeReproduction(enums.PetitHerbe), // ModeReproduction
		)
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.ID(), plant.ModeReproduction.String(), posX, posY)
	}

	// Add initial insects
	for i := 0; i < initialInsectCount; i++ {
		posX := rand.Intn(10)
		posY := rand.Intn(10)
		insect := organisme.NewInsecte(
			i,                              // ID
			0,                              // Age
			posX,                           // positionX
			posY,                           // positionY
			5,                              // Rayon
			1,                              // Vitesse
			10,                             // Energie
			10,                             // CapaciteReproduction
			5,                              // NiveauFaim
			2,                              // Hierarchie
			enums.Sexe(enums.Male),         // Sexe
			enums.MyInsect(enums.Escargot), // Espece
			1.0,                            // PeriodReproduire
			false,                          // EnvieReproduire
		)
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.ID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
	}
	return env, terr
}

//func (e *Environment) UpdateInsectPositions() {
//	for _, org := range e.Organismes {
//		if insect, ok := org.(*organisme.Insecte); ok {
//			// 更新昆虫位置的逻辑
//			newX := insect.PositionX + rand.Intn(3) - 1 // 随机 -1, 0, 1
//			newY := insect.PositionY + rand.Intn(3) - 1 // 随机 -1, 0, 1
//
//			// 确保新位置在环境范围内
//			newX = max(0, min(newX, e.Width-1))
//			newY = max(0, min(newY, e.Height-1))
//
//			insect.PositionX = newX
//			insect.PositionY = newY
//		}
//	}
//}
//
//func min(a, b int) int {
//	if a < b {
//		return a
//	}
//	return b
//}
//
//func max(a, b int) int {
//	if a > b {
//		return a
//	}
//	return b
//}
