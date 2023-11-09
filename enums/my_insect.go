package enums

type MyInsect int

const (
	Escargot MyInsect = iota
	Grillons
	Lombric
	PetitSerpent
	AraignéeSauteuse
)

// Return the string representation of the MyPlant
func (i MyInsect) String() string {
	return [...]string{"Escargot", "Grillons", "Lombric", "PetitSerpent", "AraignéeSauteuse"}[i]
}
