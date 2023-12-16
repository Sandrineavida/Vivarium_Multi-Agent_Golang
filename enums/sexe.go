package enums

// Sexe represents the sex of an organism.
type Sexe int

// Enumeration of Sexe
const (
	Male    Sexe = iota // iota is reset to 0
	Femelle             // iota increments automatically
	Hermaphrodite
)

var StringToSexe = map[string]Sexe{
	"Male":          Male,
	"Femelle":       Femelle,
	"Hermaphrodite": Hermaphrodite, //escorgot
}

// String returns the string representation of the Sexe
func (s Sexe) String() string {
	return [...]string{"Male", "Femelle", "Hermaphrodite"}[s]
}
