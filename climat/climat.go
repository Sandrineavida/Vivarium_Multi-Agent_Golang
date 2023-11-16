package climat

import "vivarium/enums"

// Climat represents the climate conditions in the simulation.
type Climat struct {
	Luminaire   int     // en lumens
	Temperature int     // en degrés celcius
	Humidite    float64 // en %
	Co2         float64 // en ppm partie par million
	O2          float64 // La mesure est exprimée en pourcentage du volume d'air (%vol)
}

// NewClimat creates a new instance of Climat with default values.qui représente un climat médian
func NewClimat() *Climat {
	return &Climat{
		Luminaire:   800,
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
	case enums.Brouillard:
		// Change climate conditions for fog
		// ... handle other cases
	case enums.SaisonSeche: // Dry Season
	case enums.Incendie: // Fire
	case enums.Tonnerre:
	}
}
