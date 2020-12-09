package utils

import (
	"net"
)

type Node struct {
	Name           string
	SubNetFileName string
	Listener       net.Listener
	// Lefs 		centralsim.Lefs
	// ClkLocal 	centralsim.TypeClock
	Partners   map[string]Partner // V t € output transition, E partner
	PetriDistr map[int]string     // For each output transition index there is an name of the target transition
	LookAhead  int
}
