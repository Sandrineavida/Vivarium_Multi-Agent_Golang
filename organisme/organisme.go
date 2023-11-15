package organisme

import (
	"vivarium/enums"
	"vivarium/terrain"
)

// Organisme defines the interface that all organisms must implement.
type Organisme interface {
	SeDeplacer(t *terrain.Terrain)
	Vieillir()
	Mourir(t *terrain.Terrain)
	//================================================
	GetID() int
	GetAge() int
	GetPosX() int
	GetPosY() int
	GetRayon() int
	GetEspece() enums.MyEspece
}

// BaseOrganisme provides the base implementation of the Organisme interface.
type BaseOrganisme struct {
	OrganismeID int
	Age         int
	PositionX   int
	PositionY   int
	Rayon       int
	Espece      enums.MyEspece
}

// NewBaseOrganisme creates a new BaseOrganisme instance.
func NewBaseOrganisme(id int, age, posX, posY, rayon int, espece enums.MyEspece) *BaseOrganisme {
	return &BaseOrganisme{
		OrganismeID: id,
		Age:         age,
		PositionX:   posX,
		PositionY:   posY,
		Rayon:       rayon,
		Espece:      espece,
	}
}

// Implementations of the Organisme interface's methods follow:

// ID returns the organism's ID.
func (bo *BaseOrganisme) ID() int {
	return bo.OrganismeID
}

// Vieillir simulates the organism aging.
func (bo *BaseOrganisme) Vieillir() {
	// Implementation of aging
	bo.Age++
}

// Mourir simulates the organism dying.
func (bo *BaseOrganisme) Mourir(t *terrain.Terrain) {
	// Implementation of dying
	// I'm not quite sure if we should clear the memory of the organism (maybe we should)
	/*...*/
	// Remove the organism from the Terrain
	t.RemoveOrganism(bo.OrganismeID, bo.PositionX, bo.PositionY)
}

func (bo *BaseOrganisme) GetID() int {
	return bo.OrganismeID
}
func (bo *BaseOrganisme) GetAge() int {
	return bo.Age
}
func (bo *BaseOrganisme) GetPosX() int {
	return bo.PositionX
}
func (bo *BaseOrganisme) GetPosY() int {
	return bo.PositionY
}
func (bo *BaseOrganisme) GetRayon() int {
	return bo.Rayon
}
func (bo *BaseOrganisme) GetEspece() enums.MyEspece {
	return bo.Espece
}
