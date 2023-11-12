package environnement

import (
	"vivarium/climat"
	"vivarium/organisme"
)

// Environment represents the simulation environment.
type Environment struct {
	Climat     *climat.Climat
	QualiteSol int
	Width      int
	Height     int
	NbPierre   int
	Organismes []*organisme.Organisme
}

// NewEnvironment creates a new instance of Environment with default values.
func NewEnvironment(width, height int) *Environment {
	return &Environment{
		Climat:     climat.NewClimat(),
		Width:      width,
		Height:     height,
		Organismes: make([]*organisme.Organisme, 0),
		// Set other attributes...
	}
}

// Simuler simulates the environment for a time step.
func (e *Environment) Simuler() {
	// Implementation of simulation step
}

// AjouterOrganisme adds a new organism to the environment.
func (e *Environment) AjouterOrganisme(o *organisme.Organisme) {
	// Implementation to add a new organism
	e.Organismes = append(e.Organismes, o)
}

// RetirerOrganisme removes an organism from the environment.
func (e *Environment) RetirerOrganisme(o *organisme.Organisme) {
	// Implementation to remove an organism
	// This might involve searching for the organism in the list and removing it.
}

// MiseAJour updates the environment state.
func (e *Environment) MiseAJour() {
	// Implementation to update the environment
	// This might involve updating the state of each organism, climate changes, etc.
}
