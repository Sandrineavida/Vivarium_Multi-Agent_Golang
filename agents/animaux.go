package agents

// fonction combat
import (
	"fmt"
	"math/rand"
	"time"
)

// Enums
type Sexe int

const (
	Male Sexe = iota
	Femelle
)

type Nourriture int

const (
	PetiteHerbe Nourriture = iota
	GrandeHerbe
	Champignon
	// autres types...
)

type ModeReproduction int

const (
	Graines ModeReproduction = iota
	Spores
	// autres modes...
)

type Meteo int

const (
	Pluie Meteo = iota
	Brouillard
	SaisonSeche
	Incendie
	Tonnerre
	// autres conditions météorologiques...
)

// Organisme
//type Organisme struct {
//	nom       string
//	age       int
//	positionX int
//	positionY int
//	rayon     int
//}

// Insecte
type Insecte struct {
	nom                     string
	age                     int
	positionX               int
	positionY               int
	rayon                   int
	sexe                    Sexe
	vitesse                 int
	sourceNourriture        Nourriture
	energie                 int
	capaciteReproduction    int
	niveauFaim              int
	periodReproduire        float64
	envieReproduire         bool
	listePourManger_Insecte []*Insecte
	listePourManger_Plante  []*Plante
	hierarchie              int
}

// Plante
type Plante struct {
	nom                 string
	age                 int
	positionX           int
	positionY           int
	rayon               int
	vitesseDeCroissance int
	etatSante           int
	modeReproduction    ModeReproduction
	adaptabilite        int
}

// Climat
type Climat struct {
	luminaire   int
	temperature int
	humidite    float64
	co2         float64
	o2          float64
}

// Environment
type Environment struct {
	climat             *Climat
	qualiteSol         int
	width              int
	height             int
	nbPierre           int
	organismes_insecte []*Insecte
	organismes_plante  []*Plante
}

func combat(insecte1, insecte2 Insecte) string {
	rand.Seed(time.Now().UnixNano())

	score1 := 0.5*insecte1.energie + 0.5*insecte1.hierarchie + 0.5*insecte1.age + 0.5*rand.Int()
	score2 := 0.5*insecte1.energie + 0.5*insecte1.hierarchie + 0.5*insecte1.age + 0.5*rand.Int()

	if score1 > score2 {
		return "Le premier animal a gagné le combat!"
	} else if score2 > score1 {
		return "Le deuxième animal a gagné le combat!"
	}
	return "Le combat s'est terminé par une égalité!"
}

func main() {
	insecte1 := Insecte{energie: 100, hierarchie: 3, age: 5}
	insecte2 := Insecte{energie: 90, hierarchie: 2, age: 3}

	result := combat(insecte1, insecte2)
	fmt.Println(result)
}
