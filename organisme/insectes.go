package organisme

import (
	"vivarium/enums"
	"vivarium/environnement"
)

// Insecte represents an insect and embeds BaseOrganisme to inherit its properties.
type Insecte struct {
	*BaseOrganisme
	Sexe                 enums.Sexe
	Vitesse              int
	Energie              int
	CapaciteReproduction int
	NiveauFaim           int
	PeriodReproduire     float64
	EnvieReproduire      bool
	ListePourManger      []string
	Hierarchie           int
}

// NewInsecte creates a new Insecte with the given attributes.
func NewInsecte(organismeID int, nom string, age, posX, posY, rayon, vitesse, energie, capaciteReproduction, niveauFaim, hierarchie int, sexe enums.Sexe, periodReproduire float64, envieReproduire bool) *Insecte {
	return &Insecte{
		BaseOrganisme:        NewBaseOrganisme(organismeID, nom, age, posX, posY, rayon),
		Sexe:                 sexe,
		Vitesse:              vitesse,
		Energie:              energie,
		CapaciteReproduction: capaciteReproduction,
		NiveauFaim:           niveauFaim,
		PeriodReproduire:     periodReproduire,
		EnvieReproduire:      envieReproduire,
		ListePourManger:      []string{},
		Hierarchie:           hierarchie,
	}
}

// Other methods (Manger, SeBattre, SeReproduire, Deplacer) need to be implemented here.

// Deplacer would need to interact with Terrain to update the organism's position.
func (i *Insecte) SeDeplacer(terrain *environnement.Terrain, positionX, positionY int) {
	// Implementation of moving behavior specific to Insecte
	// This should also update the Terrain's representation of the organism's position.
}
