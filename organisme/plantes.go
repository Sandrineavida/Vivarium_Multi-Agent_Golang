package organisme

import (
	"fmt"
	"vivarium/climat"
	"vivarium/enums"
	"vivarium/terrain"
	"vivarium/utils"
)

// Plante represents a plant and embeds BaseOrganisme to inherit its properties.
type Plante struct {
	*BaseOrganisme
	EtatSante            int
	ModeReproduction     enums.ModeReproduction
	AgeGaveBirthLastTime int
	PeriodReproduire     int
}

// NewPlante creates a new Plante with the given attributes.
func NewPlante(id, age, posX, posY, etatSante int, espece enums.MyEspece) *Plante {

	attributes := enums.SpeciesAttributes[espece]
	attributesPlante := enums.PlantAttributesMap[espece]

	return &Plante{
		BaseOrganisme: NewBaseOrganisme(id, age, posX, posY, attributesPlante.Rayon, espece,
			attributes.AgeRate, attributes.MaxAge, attributes.GrownUpAge, attributes.TooOldToReproduceAge, attributes.NbProgeniture),
		// EtatSante:            etatSante,
		EtatSante:            attributes.NiveauEnergie,
		ModeReproduction:     attributesPlante.ModeReproduction,
		AgeGaveBirthLastTime: 0,
		PeriodReproduire:     attributesPlante.PeriodReproduire,
	}
}

func (pl *Plante) CheckEtat(t *terrain.Terrain) Organisme {
	// 如果EtatSante为0，就死亡
	if pl.EtatSante <= 0 {
		pl.Mourir(t)
		return pl
	}
	return nil
}

// ========================================== MisaAJour_EtatSante ==========================================

func CanPhotosynthesize(climat climat.Climat) bool {
	return climat.Luminaire >= 20 && // 至少20%的光照
		climat.Temperature >= 10 && climat.Temperature <= 35 && // 温度在10°C至35°C之间
		climat.Humidite >= 50 && climat.Humidite <= 70 && // 湿度在50%至70%之间
		climat.Co2 >= 10 && // 至少10%的二氧化碳浓度
		climat.O2 <= 30 // 氧气浓度不超过30%
}

func IsHarshEnvironment(climat climat.Climat) bool {
	// 温度极端
	if climat.Temperature < 0 || climat.Temperature > 45 {
		return true
	}

	// 湿度极端
	if climat.Humidite < 10 || climat.Humidite > 90 {
		return true
	}

	// 光照极低
	if climat.Luminaire < 10 {
		return true
	}

	// 二氧化碳水平极端
	if climat.Co2 < 1 || climat.Co2 > 100 {
		return true
	}

	// 氧气水平极端
	if climat.O2 < 10 || climat.O2 > 30 {
		return true
	}

	return false
}

func (p *Plante) MisaAJour_EtatSante(climat climat.Climat) {
	// 先看能否进行光合作用
	if CanPhotosynthesize(climat) {
		// 如果能，就EtatSante+1(EtatSante最大为10)
		attr := enums.SpeciesAttributes[p.Espece]
		p.EtatSante = utils.Intmin(p.EtatSante+1, attr.NiveauEnergie)
		return
	} else {
		// 如果不能
		// 判断环境是否恶劣
		if IsHarshEnvironment(climat) {
			// 如果恶劣，EtatSante-1(EtatSante最小为0)
			p.EtatSante = utils.Intmax(p.EtatSante-1, 0)
			return
		}
		// 如果不恶劣，EtatSante不变
	}
}

// ========================================== End MisaAJour_EtatSante ==========================================

// ========================================== Reproduire ==========================================

// 暂时不用管Busy

func (p *Plante) CanReproduire() bool {
	//fmt.Println("操，植物可不可以生啊!!!!!!!!!!!!!!!!!!!!!", p.Age-p.AgeGaveBirthLastTime >= p.PeriodReproduire && p.Age >= p.GrownUpAge && p.Age <= p.TooOldToReproduceAge && p.EtatSante >= 5)
	return p.Age-p.AgeGaveBirthLastTime >= p.PeriodReproduire && p.Age >= p.GrownUpAge && p.Age <= p.TooOldToReproduceAge && p.EtatSante >= 5
}

func (p *Plante) Reproduire(organismes []Organisme, t *terrain.Terrain) (int, []Organisme) {
	if p.CanReproduire() {
		var sliceNewBorn []Organisme
		for i := 0; i < p.NbProgeniture; i++ {
			// posX, posY := utils.RandomPositionInRadius(p.PositionX, p.PositionY, p.Rayon)

			// // 确保坐标在 Terrain 的边界内
			// posX = utils.Intmax(utils.Intmin(posX, t.Width-1), 0)
			// posY = utils.Intmax(utils.Intmin(posY, t.Length-1), 0)
			posX, posY := utils.RandomPositionInRectangle(p.PositionX, p.PositionY, p.Rayon, 0, t.Width-1, 0, t.Length-1)

			fmt.Println("植物[", p.OrganismeID, "]生的位置: (", posX, ",", posY, ")")

			newBorn := NewPlante(-1, 0, posX, posY, 10, p.Espece)
			sliceNewBorn = append(sliceNewBorn, newBorn)
		}
		p.AgeGaveBirthLastTime = p.Age
		// fmt.Println("植物也生了!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", sliceNewBorn)
		return p.NbProgeniture, sliceNewBorn
	}
	return 0, nil
}

// ========================================== End Reproduire ==========================================

// func (p *Plante) InteragirInsecte(insecte *Insecte) {
// 	// Implementation of interaction with an insect
// }
