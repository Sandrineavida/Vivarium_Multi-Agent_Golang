package climat

import "vivarium/enums"

// Climat represents the climate conditions in the simulation.
type Climat struct {
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
	}
}
