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
	NiveauEnergie        int
}{
	PetitHerbe:       {AgeRate: 1, MaxAge: 30, GrownUpAge: 10, TooOldToReproduceAge: 30, NbProgeniture: 1, NiveauEnergie: 10},
	GrandHerbe:       {AgeRate: 1, MaxAge: 70, GrownUpAge: 15, TooOldToReproduceAge: 40, NbProgeniture: 1, NiveauEnergie: 20},
	Champignon:       {AgeRate: 2, MaxAge: 10, GrownUpAge: 2, TooOldToReproduceAge: 10, NbProgeniture: 2, NiveauEnergie: 10},
	Escargot:         {AgeRate: 1, MaxAge: 40, GrownUpAge: 15, TooOldToReproduceAge: 30, NbProgeniture: 3, NiveauEnergie: 16},
	Grillons:         {AgeRate: 2, MaxAge: 20, GrownUpAge: 5, TooOldToReproduceAge: 15, NbProgeniture: 1, NiveauEnergie: 12},
	Lombric:          {AgeRate: 1, MaxAge: 60, GrownUpAge: 20, TooOldToReproduceAge: 30, NbProgeniture: 2, NiveauEnergie: 16},
	PetitSerpent:     {AgeRate: 1, MaxAge: 70, GrownUpAge: 20, TooOldToReproduceAge: 45, NbProgeniture: 2, NiveauEnergie: 60},
	AraignéeSauteuse: {AgeRate: 2, MaxAge: 65, GrownUpAge: 5, TooOldToReproduceAge: 10, NbProgeniture: 2, NiveauEnergie: 50},
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
	PeriodReproduire int
}

// Define the characteristics of each plant
var PlantAttributesMap = map[MyEspece]PlantAttributes{
	PetitHerbe: {Rayon: 5, PeriodReproduire: 5},
	GrandHerbe: {Rayon: 4, PeriodReproduire: 20},
	Champignon: {Rayon: 8, PeriodReproduire: 6},
}

// Structure that defines insecte properties
type InsectAttributes struct {
	Rayon            int
	PeriodReproduire int
}

// Define the characteristics of each insecte
var InsectAttributesMap = map[MyEspece]InsectAttributes{
	Escargot:         {Rayon: 2, PeriodReproduire: 6},
	Grillons:         {Rayon: 1, PeriodReproduire: 6},
	Lombric:          {Rayon: 2, PeriodReproduire: 4},
	PetitSerpent:     {Rayon: 4, PeriodReproduire: 10},
	AraignéeSauteuse: {Rayon: 3, PeriodReproduire: 6},
}

// Define the speed of each insecte
var InsectSpeeds = map[MyEspece]int{
	Escargot:         1,
	Grillons:         3,
	Lombric:          2,
	PetitSerpent:     4,
	AraignéeSauteuse: 5,
}
