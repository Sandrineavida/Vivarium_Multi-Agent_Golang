package organisme

import (
	"math/rand"
	"vivarium/enums"
	"vivarium/terrain"
	"vivarium/utils"
)

// Insecte represents an insect and embeds BaseOrganisme to inherit its properties.
type Insecte struct {
	*BaseOrganisme
	Sexe                 enums.Sexe
	Espece               enums.MyInsect
	Vitesse              int
	Energie              int
	CapaciteReproduction int
	NiveauFaim           int
	PeriodReproduire     float64
	EnvieReproduire      bool
	ListePourManger      []string
	Hierarchie           int
}

// foodMap defines what each insect species can eat.
var foodMap = map[enums.MyInsect][]string{
	enums.Escargot:         {"PetitHerbe", "GrandHerbe", "Champignon"},
	enums.Grillons:         {"Champignon"},
	enums.Lombric:          {"PetitHerbe", "GrandHerbe"},
	enums.PetitSerpent:     {"Lombric"},
	enums.AraignéeSauteuse: {"Grillons"},
}

// NewInsecte creates a new Insecte with the given attributes.
func NewInsecte(organismeID int, age, posX, posY, rayon, vitesse, energie, capaciteReproduction, niveauFaim, hierarchie int,
	sexe enums.Sexe, espece enums.MyInsect, periodReproduire float64, envieReproduire bool) *Insecte {
	insecte := &Insecte{
		BaseOrganisme:        NewBaseOrganisme(organismeID, age, posX, posY, rayon),
		Sexe:                 sexe,
		Espece:               espece,
		Vitesse:              vitesse,
		Energie:              energie,
		CapaciteReproduction: capaciteReproduction,
		NiveauFaim:           niveauFaim,
		PeriodReproduire:     periodReproduire,
		EnvieReproduire:      envieReproduire,
		ListePourManger:      foodMap[espece], // Assign the diet based on the species
		Hierarchie:           hierarchie,
	}

	return insecte
}

// Other methods (Manger, SeBattre, SeReproduire, SeDeplacer) need to be implemented here.

// SeDeplacer updates the insect's position within the terrain boundaries.
func (in *Insecte) SeDeplacer(t *terrain.Terrain) {
	// Generate random movement direction
	deltaX := rand.Intn(3) - 1 // Random int in {-1, 0, 1}
	deltaY := rand.Intn(3) - 1 // Random int in {-1, 0, 1}

	// Apply velocity and constrain the new position within the terrain boundaries
	newX := utils.Intmax(0, utils.Intmin(in.positionX+deltaX*in.Vitesse, t.Width-1))
	newY := utils.Intmax(0, utils.Intmin(in.positionY+deltaY*in.Vitesse, t.Length-1))

	// Update the insect's position in the Terrain and Insecte
	t.UpdateOrganismPosition(in.organismeID, in.Espece.String(), in.positionX, in.positionY, newX, newY)
	in.positionX = newX
	in.positionY = newY
}

func (in *Insecte) Mourir() {
	// Implementation of dying
	// Might involve removing the plant from the environment.
}
