package environnement

import (
	"math/rand"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/organisme"
	"vivarium/terrain"
)

const (
	width  = 15
	height = 15
)

// Environment represents the simulation environment.
type Environment struct {
	Climat     *climat.Climat
	QualiteSol int
	Width      int
	Height     int
	NbPierre   int
	Engrais    int
	Hour       int
	Organismes []organisme.Organisme
}

// NewEnvironment creates a new instance of Environment with default values.
func NewEnvironment(width, height int) *Environment {
	return &Environment{
		Climat:     climat.NewClimat(),
		Width:      width,
		Height:     height,
		Organismes: make([]organisme.Organisme, 0),
		Hour:       0,
	}
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

// Initial number of assumptions
const (
	initPetitHerbeCount       = 100 //35
	initGrandHerbeCount       = 6
	initChampignonCount       = 4 //8
	initEscargotCount         = 20
	initGrillonsCount         = 4
	initLombricCount          = 5
	initPetitSerpentCount     = 2
	initAraignéeSauteuseCount = 4
)

var Insects []*organisme.Insecte
var Plants []*organisme.Plante

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
	env := NewEnvironment(width, height)
	terr := terrain.NewTerrain(width, height)
	terr.Meteo = env.Climat.Meteo
	terr.Luminaire = env.Climat.Luminaire
	terr.Temperature = env.Climat.Temperature
	terr.Humidite = env.Climat.Humidite
	terr.O2 = env.Climat.O2
	terr.Co2 = env.Climat.Co2

	// Add initial plants
	// func NewPlante(id, age, posX, posY, rayon, vitesseDeCroissance, etatSante, adaptabilite int, modeReproduction enums.ModeReproduction, espece enums.MyEspece)
	// PetitHerbe
	for i := 0; i < initPetitHerbeCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		plant := organisme.NewPlante(
			id,   // ID
			0,    // Age
			posX, // positionX
			posY, // positionY
			enums.PetitHerbe,
		)
		id = id + 1
		//env.AjouterOrganisme(toOrganisme(plant))
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.GetID(), plant.Espece.String(), posX, posY)
		Plants = append(Plants, plant)
	}
	// GrandHerbe
	for i := 0; i < initGrandHerbeCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		plant := organisme.NewPlante(
			id,   // ID
			0,    // Age
			posX, // positionX
			posY, // positionY
			enums.GrandHerbe,
		)
		//env.AjouterOrganisme(toOrganisme(plant))
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.GetID(), plant.Espece.String(), posX, posY)
		Plants = append(Plants, plant)
		id = id + 1
	}
	// Champignon
	for i := 0; i < initChampignonCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		plant := organisme.NewPlante(
			id,   // ID
			0,    // Age
			posX, // positionX
			posY, // positionY
			enums.Champignon,
		)
		//env.AjouterOrganisme(toOrganisme(plant))
		env.AjouterOrganisme(plant)
		terr.AddOrganism(plant.GetID(), plant.Espece.String(), posX, posY)
		Plants = append(Plants, plant)
		id = id + 1
	}

	// Add initial insects
	// func NewInsecte(organismeID int, age, posX, posY, rayon, vitesse, energie, capaciteReproduction, niveauFaim int,
	//	sexe enums.Sexe, espece enums.MyEspece, periodReproduire time.Duration, envieReproduire bool)

	// Escargot - Hermaphrodite
	for i := 0; i < initEscargotCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                              // ID
			0,                               // Age
			posX,                            // positionX
			posY,                            // positionY
			enums.Sexe(enums.Hermaphrodite), // Sexe
			enums.Escargot,                  // espace
			false,                           // EnvieReproduire
		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}
	// Grillons - Male
	for i := 0; i < initGrillonsCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                     // ID
			0,                      // Age
			posX,                   // positionX
			posY,                   // positionY
			enums.Sexe(enums.Male), // Sexe
			enums.Grillons,         // espace
			false,                  // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}
	// Grillons - Femelle
	for i := 0; i < initGrillonsCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                        // ID
			0,                         // Age
			posX,                      // positionX
			posY,                      // positionY
			enums.Sexe(enums.Femelle), // Sexe
			enums.Grillons,            // espace
			false,                     // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}

	// Lombric - Hermaphrodite
	for i := 0; i < initLombricCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                              // ID
			0,                               // Age
			posX,                            // positionX
			posY,                            // positionY
			enums.Sexe(enums.Hermaphrodite), // Sexe
			enums.Lombric,                   // espace
			false,                           // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}

	// AraignéeSauteuse - Male
	for i := 0; i < initAraignéeSauteuseCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                     // ID
			0,                      // Age
			posX,                   // positionX
			posY,                   // positionY
			enums.Sexe(enums.Male), // Sexe
			enums.AraignéeSauteuse, // espace
			false,                  // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}
	// AraignéeSauteuse - Femelle
	for i := 0; i < initAraignéeSauteuseCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                        // ID
			0,                         // Age
			posX,                      // positionX
			posY,                      // positionY
			enums.Sexe(enums.Femelle), // Sexe
			enums.AraignéeSauteuse,    // espace
			false,                     // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}

	// PetitSerpent - Male
	for i := 0; i < initPetitSerpentCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                     // ID
			0,                      // Age
			posX,                   // positionX
			posY,                   // positionY
			enums.Sexe(enums.Male), // Sexe
			enums.PetitSerpent,     // espace
			false,                  // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}

	for i := 0; i < initPetitSerpentCount; i++ {
		posX := rand.Intn(env.Width)
		posY := rand.Intn(env.Height)
		insect := organisme.NewInsecte(
			id,                        // ID
			0,                         // Age
			posX,                      // positionX
			posY,                      // positionY
			enums.Sexe(enums.Femelle), // Sexe
			enums.PetitSerpent,        // espace
			false,                     // EnvieReproduire

		)
		//env.AjouterOrganisme(toOrganisme(insect))
		env.AjouterOrganisme(insect)
		terr.AddOrganism(insect.GetID(), insect.Espece.String(), posX, posY)
		Insects = append(Insects, insect) // Used to provide to the main function to allow all insects to move randomly
		id = id + 1
	}

	return env, terr, id
}
