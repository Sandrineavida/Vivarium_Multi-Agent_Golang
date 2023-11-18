package organisme

import (
	"vivarium/enums"
	"vivarium/terrain"
)

// Plante represents a plant and embeds BaseOrganisme to inherit its properties.
type Plante struct {
	*BaseOrganisme
	VitesseDeCroissance  int
	EtatSante            int
	ModeReproduction     enums.ModeReproduction
	Adaptabilite         int
	AgeGaveBirthLastTime int
}

// NewPlante creates a new Plante with the given attributes.
func NewPlante(id, age, posX, posY, etatSante, adaptabilite int, espece enums.MyEspece) *Plante {

	attributes := enums.SpeciesAttributes[espece]
	attributesPlante := enums.PlantAttributesMap[espece]

	return &Plante{
		BaseOrganisme: NewBaseOrganisme(id, age, posX, posY, attributesPlante.Rayon, espece,
			attributes.AgeRate, attributes.MaxAge, attributes.GrownUpAge, attributes.NbProgeniture, attributes.TooOldToReproduceAge),
		VitesseDeCroissance:  attributesPlante.VitesseDeCroissance,
		EtatSante:            etatSante,
		ModeReproduction:     attributesPlante.ModeReproduction,
		Adaptabilite:         adaptabilite,
		AgeGaveBirthLastTime: 0,
	}
}

// Implement the methods specific to Plante here.

func (p *Plante) Photosynthese() {
	// Implementation of photosynthesis behavior
}

func (p *Plante) Grandir() {
	// Implementation of growing behavior
}

func (p *Plante) Reproduire() *Plante {
	// Implementation of reproduction behavior
	return nil // Replace with actual implementation
}

func (p *Plante) InteragirInsecte(insecte *Insecte) {
	// Implementation of interaction with an insect
}

// Implement the Organisme interface methods (SeDeplacer, Vieillir, Mourir).

func (p *Plante) SeDeplacer(t *terrain.Terrain) {
	// Plants might not move, so this could be a no-op or handled differently.
}

// func (p *Plante) Vieillir() {
// 	// Implementation of aging
// 	p.Age++
// }

// func (p *Plante) Mourir(t *terrain.Terrain) {
// 	// Implementation of dying
// 	// Might involve removing the plant from the environment.
// }
