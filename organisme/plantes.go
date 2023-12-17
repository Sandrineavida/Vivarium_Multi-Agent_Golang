package organisme

import (
	"time"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/terrain"
	"vivarium/utils"
)

// Plante represents a plant and embeds BaseOrganisme to inherit its properties.
type Plante struct {
	*BaseOrganisme
	EtatSante            int
	AgeGaveBirthLastTime int
	PeriodReproduire     int
	NbParts              int  // Only GrandHerbe has this property
	IsBeingEaten         bool // Only GrandHerbe has this property
	// For animation (Ebiten)
	IsReproduire bool
	IsNormal     bool
}

// NewPlante creates a new Plante with the given attributes.
func NewPlante(id, age, posX, posY int, espece enums.MyEspece) *Plante {

	attributes := enums.SpeciesAttributes[espece]
	attributesPlante := enums.PlantAttributesMap[espece]

	return &Plante{
		BaseOrganisme: NewBaseOrganisme(id, age, posX, posY, attributesPlante.Rayon, espece,
			attributes.AgeRate, attributes.MaxAge, attributes.GrownUpAge, attributes.TooOldToReproduceAge, attributes.NbProgeniture, false),
		EtatSante:            attributes.NiveauEnergie,
		AgeGaveBirthLastTime: 0,
		PeriodReproduire:     attributesPlante.PeriodReproduire,
		NbParts:              4,
		IsBeingEaten:         false,
		IsReproduire:         false,
		IsNormal:             true,
	}
}

func (pl *Plante) CheckEtat(t *terrain.Terrain) Organisme {
	// if EtatSante <= 0, then mourir
	if pl.EtatSante <= 0 {
		pl.Mourir(t)
		return pl
	}
	return nil
}

// ========================================== MisaAJour_EtatSante ==========================================

func CanPhotosynthesize(climat climat.Climat) bool {
	return climat.Luminaire >= 20 &&
		climat.Temperature >= 10 && climat.Temperature <= 35 &&
		climat.Humidite >= 50 && climat.Humidite <= 70 &&
		climat.Co2 >= 10 &&
		climat.O2 <= 30
}

func DegreeHarshEnv(climat climat.Climat) int {
	degree := 0 // degree = [0, 17]
	// Extreme Temperature
	if climat.Temperature < 0 {
		degree = 1
	}
	if climat.Temperature > 40 && climat.Temperature <= 55 {
		degree = 2
	}
	if climat.Temperature > 55 {
		degree = 15
	}

	// Extreme Co2
	if climat.Co2 < 1 || climat.Co2 > 100 {
		degree += 1
	}

	// Extreme O2
	if climat.O2 < 10 || climat.O2 > 30 {
		degree += 1
	}

	return degree
}

func (p *Plante) MisaAJour_EtatSante(climat climat.Climat) {
	// see if it can photosynthesize
	if CanPhotosynthesize(climat) {
		// if can, then EtatSante+1
		if p.Espece != enums.Champignon { // Champignon can't photosynthesize
			attr := enums.SpeciesAttributes[p.Espece]
			p.EtatSante = utils.Intmin(p.EtatSante+1, attr.NiveauEnergie)
		}
		return
	} else {
		// Decrease EtatSante according to the degree of harsh environment
		harshenv_degree := DegreeHarshEnv(climat)
		p.EtatSante = utils.Intmax(p.EtatSante-harshenv_degree, 0)
		return
		// if the environment is not harsh, EtatSante won't get changed
	}
}

// ========================================== End MisaAJour_EtatSante ==========================================

// ========================================== Reproduire ==========================================

func (p *Plante) CanReproduire() bool {
	return p.Age-p.AgeGaveBirthLastTime >= p.PeriodReproduire && p.Age >= p.GrownUpAge && p.Age <= p.TooOldToReproduceAge && p.EtatSante >= 5
}

func (p *Plante) Reproduire(organismes []Organisme, t *terrain.Terrain) (int, []Organisme) {
	if p.CanReproduire() {
		var sliceNewBorn []Organisme
		for i := 0; i < p.NbProgeniture; i++ {
			// randomly pick up a position within the radius of the parent plant
			posX, posY := utils.RandomPositionInRectangle(p.PositionX, p.PositionY, p.Rayon, 0, t.Width-1, 0, t.Length-1)
			newBorn := NewPlante(-1, 0, posX, posY, p.Espece)
			sliceNewBorn = append(sliceNewBorn, newBorn)
		}
		p.AgeGaveBirthLastTime = p.Age
		return p.NbProgeniture, sliceNewBorn
	}
	p.IsNormal = false
	p.IsReproduire = true
	defer func() {
		time.Sleep(2 * timeSleep * time.Millisecond)
		p.IsReproduire = false
		p.IsNormal = true
	}()
	return 0, nil
}

// ========================================== End Reproduire ==========================================
