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
