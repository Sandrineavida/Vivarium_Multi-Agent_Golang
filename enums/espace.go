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
	PetitSerpent:     {AgeRate: 1, MaxAge: 70, GrownUpAge: 20, TooOldToReproduceAge: 45, NbProgeniture: 2, NiveauEnergie: 30},
	AraignéeSauteuse: {AgeRate: 3, MaxAge: 65, GrownUpAge: 5, TooOldToReproduceAge: 10, NbProgeniture: 2, NiveauEnergie: 36},
}

var InsectSpeeds = map[MyEspece]int{
	Escargot:         2,  // 蜗牛速度较慢
	Grillons:         6,  // 蟋蟀速度中等
	Lombric:          3,  // 蚯蚓速度较慢
	PetitSerpent:     8,  // 小蛇速度较快
	AraignéeSauteuse: 10, // 跳蛛速度非常快
	// ...其他昆虫物种
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
	PetitHerbe: {Rayon: 5, ModeReproduction: Graine, PeriodReproduire: 5},
	GrandHerbe: {Rayon: 4, ModeReproduction: Graine, PeriodReproduire: 20},
	Champignon: {Rayon: 8, ModeReproduction: Spore, PeriodReproduire: 6},
}

// Structure that defines insecte properties
type InsectAttributes struct {
	Rayon            int
	PeriodReproduire int
}

// Define the characteristics of each insecte
var InsectAttributesMap = map[MyEspece]InsectAttributes{
	Escargot:         {Rayon: 2, PeriodReproduire: 6},  //理论上可以繁殖 次
	Grillons:         {Rayon: 1, PeriodReproduire: 6},  //理论上可以繁殖 次
	Lombric:          {Rayon: 2, PeriodReproduire: 4},  //理论上可以繁殖 次
	PetitSerpent:     {Rayon: 4, PeriodReproduire: 10}, //理论上可以繁殖 次
	AraignéeSauteuse: {Rayon: 3, PeriodReproduire: 6},  //理论上可以繁殖 次
}
