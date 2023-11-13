package organisme

import (
	"vivarium/terrain"
)

// Organisme defines the interface that all organisms must implement.
type Organisme interface {
	Vieillir()
	Mourir()
	ID() int
}

// BaseOrganisme provides the base implementation of the Organisme interface.
type BaseOrganisme struct {
	organismeID int
	age         int
	positionX   int
	positionY   int
	rayon       int
}

// NewBaseOrganisme creates a new BaseOrganisme instance.
func NewBaseOrganisme(id int, age, posX, posY, rayon int) *BaseOrganisme {
	return &BaseOrganisme{
		organismeID: id,
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

// Vieillir simulates the organism aging.
func (bo *BaseOrganisme) Vieillir() {
	// Implementation of aging
	bo.age++
}

// Mourir simulates the organism dying.
func (bo *BaseOrganisme) Mourir(t *terrain.Terrain) {
	// Implementation of dying
	// I'm not quite sure if we should clear the memory of the organism (maybe we should)
	/*...*/
	// Remove the organism from the Terrain
	t.RemoveOrganism(bo.organismeID, bo.positionX, bo.positionY)
}

// Additional methods or attributes should be added as necessary, based on the specific logic and requirements
// of your simulation.
