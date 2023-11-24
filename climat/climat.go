package climat

import (
	"vivarium/enums"
	"vivarium/utils"
)

// Climat represents the climate conditions in the simulation.
type Climat struct {
	Meteo       enums.Meteo
	Luminaire   int     //0-100 %
	Temperature int     //-5-400 ℃
	Humidite    float32 //0-100 %
	Co2         float32 //0-100 %
	O2          float32 //0-100 %
}

// NewClimat creates a new instance of Climat with default values.
func NewClimat() *Climat {
	return &Climat{
		// Initialize with default values or parameters
		Meteo:       enums.Rien,
		Luminaire:   50,    // 默认光照50%
		Temperature: 20,    // 默认温度20°C
		Humidite:    50.0,  // 默认湿度50%
		Co2:         50.0,  // 默认二氧化碳50%
		O2:          20.95, // 默认氧气20.95%（大气中的平均水平）
	}
}

// ChangerConditions changes the current weather conditions based on Meteo.
func (c *Climat) ChangerConditions(meteo enums.Meteo) (engrais int) {
	// Implementation of how different weather conditions affect the climate
	switch meteo {
	case enums.Pluie:
		// Change climate conditions for rain
		c.Meteo = enums.Pluie
		c.Humidite = utils.Float32min(c.Humidite+10.5, 100)
		c.Temperature = utils.Intmax(c.Temperature-2, -5)
		c.O2 = utils.Float32max(c.O2+2.5, 0)
	case enums.Brouillard:
		// Change climate conditions for fog
		c.Meteo = enums.Brouillard
		c.Humidite = utils.Float32min(c.Humidite+5.5, 100)
		c.Temperature = utils.Intmax(c.Temperature-1, -5)
		c.O2 = utils.Float32max(c.O2+1.5, 0)
	case enums.SaisonSeche:
		// Change climate conditions for dry season
		c.Meteo = enums.SaisonSeche
		c.Humidite = utils.Float32max(c.Humidite-3.5, 0)
		c.Temperature = utils.Intmin(c.Temperature+1, 40)
		c.Co2 = utils.Float32min(c.Co2+2.5, 0)
	case enums.Incendie:
		// Change climate conditions for fire
		c.Meteo = enums.Incendie
		c.Humidite = utils.Float32max(c.Humidite-27.5, 0)
		c.Temperature = utils.Intmax(c.Temperature+200, 400)
		c.Co2 = utils.Float32max(c.Co2+32.5, 100)
		c.O2 = utils.Float32min(c.O2-32.5, 0)
	case enums.Tonnerre:
		// Change climate conditions for thunder
		c.Meteo = enums.Tonnerre
		engrais = 20
	case enums.Rien:
		c.Meteo = enums.Rien
		c.Humidite = 50.0 // 默认湿度50%
		c.Co2 = 50.0      // 默认二氧化碳50%
		c.O2 = 20.95
	}
	return engrais
}

func (c *Climat) UpdateClimat_24H(hour int) {
	// 只更改光照和温度
	// 0 6 12 18
	switch hour {
	case 0:
		c.Luminaire = 10
		c.Temperature = 5
	case 6:
		c.Luminaire = 40
		c.Temperature = 10
	case 12:
		c.Luminaire = 100
		c.Temperature = 20
	case 18:
		c.Luminaire = 50
		c.Temperature = 9
	}
}

// func CanPhotosynthesize(climat climat.Climat) bool {
// 	return climat.Luminaire >= 20 && // 至少20%的光照
// 		climat.Temperature >= 10 && climat.Temperature <= 35 && // 温度在10°C至35°C之间
// 		climat.Humidite >= 50 && climat.Humidite <= 70 && // 湿度在50%至70%之间
// 		climat.Co2 >= 10 && // 至少10%的二氧化碳浓度
// 		climat.O2 <= 30 // 氧气浓度不超过30%
// }
