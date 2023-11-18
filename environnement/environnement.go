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
func (e *Environment) RetirerOrganisme(o organisme.Organisme) {
	for i, org := range e.Organismes {
		if org.GetID() == o.GetID() {
			e.Organismes = append(e.Organismes[:i], e.Organismes[i+1:]...)
			break
		}
	}
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
	initialInsectCount = 10
)

var Insects []*organisme.Insecte

type OrganismeInfo struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Species    string `json:"species"`
	Position_X int    `json:"position_x"`
	Position_Y int    `json:"position_y"`
}

func (e *Environment) GetAllOrganisms() []organisme.Organisme {
	return e.Organismes
}

// InitializeEcosystem creates and initializes the environment and creatures of the ecosystem
func InitializeEcosystem(id int) (*Environment, *terrain.Terrain, int) {
	// Create environment instance
	env := NewEnvironment(10, 10)
	terr := terrain.NewTerrain(10, 10)

	// Add initial plants
	// func NewPlante(id, age, posX, posY, rayon, vitesseDeCroissance, etatSante, adaptabilite int, modeReproduction enums.ModeReproduction, espece enums.MyEspece)
	for i := 0; i < initialPlantCount; i++ {
		posX := rand.Intn(10)
		posY := rand.Intn(10)
		plant := organisme.NewPlante(
			id,   // ID
			0,    // Age
			posX, // positionX
			posY, // positionY
			100,  // EtatSante
			1,    // Adaptabilite
			enums.PetitHerbe,
		)
		id = id + 1
		//env.AjouterOrganisme(toOrganisme(plant))
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.GetID(), plant.Espece.String(), posX, posY)
	}
	for i := 0; i < 5; i++ {
		posX := rand.Intn(10)
		posY := rand.Intn(10)
		plant := organisme.NewPlante(
			id,   // ID
			0,    // Age
			posX, // positionX
			posY, // positionY
			100,  // EtatSante
			1,    // Adaptabilite
			enums.GrandHerbe,
		)
		//env.AjouterOrganisme(toOrganisme(plant))
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.GetID(), plant.Espece.String(), posX, posY)
		id = id + 1
	}

	// Add initial insects
	// func NewInsecte(organismeID int, age, posX, posY, rayon, vitesse, energie, capaciteReproduction, niveauFaim int,
	//	sexe enums.Sexe, espece enums.MyEspece, periodReproduire time.Duration, envieReproduire bool)
	for i := 0; i < initialInsectCount; i++ {
		posX := rand.Intn(10)
		posY := rand.Intn(10)
		insect := organisme.NewInsecte(
			id,                     // ID
			0,                      // Age
			posX,                   // positionX
			posY,                   // positionY
			1,                      // Vitesse
			10,                     // Energie
			enums.Sexe(enums.Male), // Sexe
			enums.Escargot,         // espace
			false,                  // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}
	return env, terr, id
}
