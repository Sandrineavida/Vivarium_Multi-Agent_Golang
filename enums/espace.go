package enums

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
	AgeRate              int
	MaxAge               int
	GrownUpAge           int
	TooOldToReproduceAge int
	NbProgeniture        int
}{
	PetitHerbe:       {AgeRate: 1, MaxAge: 30, GrownUpAge: 10, TooOldToReproduceAge: 30, NbProgeniture: 1},
	GrandHerbe:       {AgeRate: 1, MaxAge: 70, GrownUpAge: 15, TooOldToReproduceAge: 40, NbProgeniture: 1},
	Champignon:       {AgeRate: 2, MaxAge: 10, GrownUpAge: 2, TooOldToReproduceAge: 10, NbProgeniture: 5},
	Escargot:         {AgeRate: 1, MaxAge: 40, GrownUpAge: 15, TooOldToReproduceAge: 25, NbProgeniture: 4},
	Grillons:         {AgeRate: 2, MaxAge: 20, GrownUpAge: 5, TooOldToReproduceAge: 15, NbProgeniture: 3},
	Lombric:          {AgeRate: 1, MaxAge: 60, GrownUpAge: 20, TooOldToReproduceAge: 30, NbProgeniture: 1},
	PetitSerpent:     {AgeRate: 1, MaxAge: 70, GrownUpAge: 20, TooOldToReproduceAge: 35, NbProgeniture: 1},
	AraignéeSauteuse: {AgeRate: 3, MaxAge: 65, GrownUpAge: 5, TooOldToReproduceAge: 10, NbProgeniture: 2},
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
	Rayon            int
	ModeReproduction ModeReproduction
	PeriodReproduire int
}

// Define the characteristics of each plant
var PlantAttributesMap = map[MyEspece]PlantAttributes{
	PetitHerbe: {Rayon: 2, ModeReproduction: Graine, PeriodReproduire: 2},
	GrandHerbe: {Rayon: 1, ModeReproduction: Graine, PeriodReproduire: 3},
	Champignon: {Rayon: 5, ModeReproduction: Spore, PeriodReproduire: 1},
}

// Structure that defines insecte properties
type InsectAttributes struct {
	Rayon            int
	PeriodReproduire int
}

// Define the characteristics of each insecte
var InsectAttributesMap = map[MyEspece]InsectAttributes{
	Escargot:         {Rayon: 2, PeriodReproduire: 3},  //理论上可以繁殖3次
	Grillons:         {Rayon: 1, PeriodReproduire: 3},  //理论上可以繁殖3次
	Lombric:          {Rayon: 2, PeriodReproduire: 2},  //理论上可以繁殖4-5次
	PetitSerpent:     {Rayon: 4, PeriodReproduire: 10}, //理论上可以繁殖1次
	AraignéeSauteuse: {Rayon: 3, PeriodReproduire: 6},  //理论上可以繁殖1次
}
