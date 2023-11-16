package climat

import "vivarium/enums"

// Climat represents the climate conditions in the simulation.
type Climat struct {
	Luminaire   int     // en lux https://fr.wikipedia.org/wiki/Lumi%C3%A8re_du_jour
	Temperature int     // en degrés celcius
	Humidite    float64 // en %
	Co2         float64 // en ppm partie par million
	O2          float64 // La mesure est exprimée en pourcentage du volume d'air (%vol)
}

// NewClimat creates a new instance of Climat with default values.qui représente un climat médian
func NewClimat() *Climat {
	return &Climat{
		Luminaire:   20000,
		Temperature: 20,
		Humidite:    50,
		Co2:         400,
		O2:          21,
	}
}

// ChangerConditions changes the current weather conditions based on Meteo.
func (c *Climat) ChangerConditions(meteo enums.Meteo) {
	// Implementation of how different weather conditions affect the climate
	switch meteo {
	case enums.Pluie:
		// Change climate conditions for rain
		c.Humidite = 100
		c.Luminaire = 100

	case enums.Brouillard:
		c.Humidite = 90
		c.Luminaire = 50
		// Change climate conditions for fog
		// ... handle other cases
	case enums.SaisonSeche:
		c.Temperature = 45
		c.Humidite = 15 // Dry Season
	case enums.Incendie: // Fire
		c.Temperature = 100
		c.Co2 = 30000
		c.O2 = 15
	case enums.Tonnerre:

	}
}
