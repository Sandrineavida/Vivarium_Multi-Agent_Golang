package organisme

import (
	"vivarium/enums"
	"vivarium/terrain"
)

// Plante represents a plant and embeds BaseOrganisme to inherit its properties.
type Plante struct {
	*BaseOrganisme
	VitesseDeCroissance int
	EtatSante           int
	ModeReproduction    enums.ModeReproduction
	Adaptabilite        int
}

// NewPlante creates a new Plante with the given attributes.
func NewPlante(id, age, posX, posY, rayon, vitesseDeCroissance, etatSante, adaptabilite int, modeReproduction enums.ModeReproduction, espece enums.MyEspece) *Plante {

	attributes := enums.SpeciesAttributes[espece]

	return &Plante{
		BaseOrganisme:       NewBaseOrganisme(id, age, posX, posY, rayon, espece, attributes.AgeRate, attributes.MaxAge),
		VitesseDeCroissance: vitesseDeCroissance,
		EtatSante:           etatSante,
		ModeReproduction:    modeReproduction,
		Adaptabilite:        adaptabilite,
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
