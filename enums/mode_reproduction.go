package enums

// ModeReproduction represents the reproduction mode of an organism.
type ModeReproduction int

// Enumeration of ModeReproduction
const (
	Vivipare ModeReproduction = iota //胎生
	Ovipare                          //卵生
	Graine                           //种子
	Spore                            //孢子
)

// String returns the string representation of the ModeReproduction
func (m ModeReproduction) String() string {
	return [...]string{"Vivipare", "Ovipare", "Graine", "Spore"}[m]
}
