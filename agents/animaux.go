package agents

// fonction combat
import (
	"fmt"
	"math/rand"
	"time"
)

type Animal struct {
	Energie    float64
	Hiérarchie float64
	Âge        float64
}

func combat(animal1, animal2 Animal) string {
	rand.Seed(time.Now().UnixNano())

	score1 := 0.5*animal1.Energie + 0.5*animal1.Hiérarchie + 0.5*animal1.Âge + 0.5*rand.Float64()
	score2 := 0.5*animal2.Energie + 0.5*animal2.Hiérarchie + 0.5*animal2.Âge + 0.5*rand.Float64()

	if score1 > score2 {
		return "Le premier animal a gagné le combat!"
	} else if score2 > score1 {
		return "Le deuxième animal a gagné le combat!"
	}
	return "Le combat s'est terminé par une égalité!"
}

func main() {
	animal1 := Animal{Energie: 100, Hiérarchie: 3, Âge: 5}
	animal2 := Animal{Energie: 90, Hiérarchie: 2, Âge: 3}

	result := combat(animal1, animal2)
	fmt.Println(result)
}
