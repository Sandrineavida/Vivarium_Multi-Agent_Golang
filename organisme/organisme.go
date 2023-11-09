package organisme

import (
	"vivarium/environnement"
)

// Organisme defines the interface that all organisms must implement.
type Organisme interface {
	SeDeplacer(t *environnement.Terrain, positionX, positionY int) // maybe only for insects; maybe move only 1 unit of distance at a time and the direction is random
	Vieillir()
	Mourir()
	ID() int
}

// BaseOrganisme provides the base implementation of the Organisme interface.
type BaseOrganisme struct {
	organismeID int
	nom         string
	age         int
	positionX   int
	positionY   int
	rayon       int
}

// NewBaseOrganisme creates a new BaseOrganisme instance.
func NewBaseOrganisme(id int, nom string, age, posX, posY, rayon int) *BaseOrganisme {
	return &BaseOrganisme{
		organismeID: id,
		nom:         nom,
		age:         age,
		positionX:   posX,
		positionY:   posY,
		rayon:       rayon,
	}
}

// Implementations of the Organisme interface's methods follow:

// ID returns the organism's ID.
func (bo *BaseOrganisme) ID() int {
	return bo.organismeID
}

// SeDeplacer updates the organism's position.
func (bo *BaseOrganisme) SeDeplacer(t *environnement.Terrain, positionX, positionY int) {
	// Update the organism's position in the Terrain
	// This method might need to access the Terrain instance to update the organism's position.
	// bo.positionX = positionX
	// bo.positionY = positionY
	// Update the Terrain with the new position
	// t.UpdateOrganismPosition(bo.organismeID, positionX, positionY, oldX, oldY)
}

// Vieillir simulates the organism aging.
func (bo *BaseOrganisme) Vieillir() {
	// Implementation of aging
	bo.age++
}

// Mourir simulates the organism dying.
func (bo *BaseOrganisme) Mourir() {
	// Implementation of dying
	// This may require interaction with the Terrain or Environment to remove the organism.
}

// Additional methods or attributes should be added as necessary, based on the specific logic and requirements
// of your simulation.
