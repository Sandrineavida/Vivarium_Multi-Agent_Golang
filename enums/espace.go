package enums

import "time"

type MyEspece int

const (
	PetitHerbe MyEspece = iota
	GrandHerbe
	Champignon
	Escargot
	Grillons
	Lombric
	PetitSerpent
	AraignéeSauteuse
)

// Return the string representation of the MyPlant
func (i MyEspece) String() string {
	return [...]string{"PetitHerbe", "GrandHerbe", "Champignon", "Escargot", "Grillons", "Lombric", "PetitSerpent", "AraignéeSauteuse"}[i]
}

var SpeciesAttributes = map[MyEspece]struct {
	AgeRate int
	MaxAge  int
}{
	PetitHerbe:       {AgeRate: 1, MaxAge: 30},
	GrandHerbe:       {AgeRate: 1, MaxAge: 40},
	Champignon:       {AgeRate: 2, MaxAge: 10},
	Escargot:         {AgeRate: 1, MaxAge: 50},
	Grillons:         {AgeRate: 2, MaxAge: 20},
	Lombric:          {AgeRate: 1, MaxAge: 60},
	PetitSerpent:     {AgeRate: 1, MaxAge: 70},
	AraignéeSauteuse: {AgeRate: 3, MaxAge: 15},
}

var StringToMyEspece = map[string]MyEspece{
	"PetitHerbe":       PetitHerbe,
	"GrandHerbe":       GrandHerbe,
	"Champignon":       Champignon,
	"Escargot":         Escargot,
	"Grillons":         Grillons,
	"Lombric":          Lombric,
	"PetitSerpent":     PetitSerpent,
	"AraignéeSauteuse": AraignéeSauteuse,
}

// Structure that defines plant properties
type PlantAttributes struct {
	Rayon               int
	VitesseDeCroissance int
	ModeReproduction    ModeReproduction
}

// Define the characteristics of each plant
var PlantAttributesMap = map[MyEspece]PlantAttributes{
	PetitHerbe: {Rayon: 2, VitesseDeCroissance: 1, ModeReproduction: Graine},
	GrandHerbe: {Rayon: 3, VitesseDeCroissance: 2, ModeReproduction: Graine},
	Champignon: {Rayon: 1, VitesseDeCroissance: 3, ModeReproduction: Spore},
}

// Structure that defines insecte properties
type InsectAttributes struct {
	Rayon            int
	PeriodReproduire time.Duration
}

// Define the characteristics of each insecte
var InsectAttributesMap = map[MyEspece]InsectAttributes{
	Escargot:         {Rayon: 2, PeriodReproduire: time.Hour * 24},
	Grillons:         {Rayon: 1, PeriodReproduire: time.Hour * 12},
	Lombric:          {Rayon: 2, PeriodReproduire: time.Hour * 48},
	PetitSerpent:     {Rayon: 4, PeriodReproduire: time.Hour * 72},
	AraignéeSauteuse: {Rayon: 3, PeriodReproduire: time.Hour * 36},
}
