package organisme

import (
	"fmt"
	"vivarium/enums"
	"vivarium/terrain"
)

// Organisme defines the interface that all organisms must implement.
type Organisme interface {
	// SeDeplacer(t *terrain.Terrain)
	Vieillir(t *terrain.Terrain)
	Mourir(t *terrain.Terrain)
	CheckEtat(t *terrain.Terrain) Organisme

	//================================================
	GetID() int
	GetAge() int
	GetPosX() int
	GetPosY() int
	GetRayon() int
	GetEspece() enums.MyEspece
	SetID(newID int)
	GetEtat() bool
}

// BaseOrganisme provides the base implementation of the Organisme interface.
type BaseOrganisme struct {
	OrganismeID          int
	Age                  int
	PositionX            int
	PositionY            int
	Rayon                int
	Espece               enums.MyEspece
	AgeRate              int
	MaxAge               int
	GrownUpAge           int
	TooOldToReproduceAge int
	NbProgeniture        int
	Busy                 bool // Lock for the organism
	IsInsecte            bool
	IsDying              bool // To be used by ebiten to render the dying animation of the insect
}

// NewBaseOrganisme creates a new BaseOrganisme instance.
func NewBaseOrganisme(id int, age, posX, posY, rayon int, espece enums.MyEspece, ageRate int, maxAge int, grownUpAge, tooOldToRep, nbProgeniture int, isInsecte bool) *BaseOrganisme {
	return &BaseOrganisme{
		OrganismeID:          id,
		Age:                  age,
		PositionX:            posX,
		PositionY:            posY,
		Rayon:                rayon,
		Espece:               espece,
		AgeRate:              ageRate,
		MaxAge:               maxAge,
		GrownUpAge:           grownUpAge,
		TooOldToReproduceAge: tooOldToRep,
		NbProgeniture:        nbProgeniture,
		Busy:                 false,
		IsInsecte:            isInsecte,
		IsDying:              false,
	}
}

// Vieillir simulates the organism aging.
func (bo *BaseOrganisme) Vieillir(t *terrain.Terrain) {
	bo.Busy = true
	defer func() { bo.Busy = false }()

	bo.Age += bo.AgeRate
	if bo.Age >= bo.MaxAge {
		// Reaching the maximum lifespan, the organism should die
		fmt.Println(bo.GetID(), "died because of old age")
		bo.Mourir(t)
	}
}

// Mourir simulates the organism dying.
func (bo *BaseOrganisme) Mourir(t *terrain.Terrain) {
	// Implementation of dying
	// Remove the organism from the Terrain
	t.RemoveOrganism(bo.OrganismeID, bo.PositionX, bo.PositionY)
	bo.IsDying = true
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
func (bo *BaseOrganisme) GetEtat() bool {
	return bo.IsDying
}
func (bo *BaseOrganisme) SetID(newID int) {
	bo.OrganismeID = newID
}
