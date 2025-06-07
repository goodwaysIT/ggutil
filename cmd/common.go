package cmd

import (
	"sync"
	// "github.com/urfave/cli/v2" // No longer needed here if all Run... functions are moved
)

// globalGGHomes stores the list of OGG Home paths set by the main package
var globalGGHomes []string
var globalGGHomesMutex sync.RWMutex

// SetGlobalGGHomes is called by the main package to set the parsed OGG Home path list
func SetGlobalGGHomes(homes []string) {
	globalGGHomesMutex.Lock()
	defer globalGGHomesMutex.Unlock()
	globalGGHomes = homes
}

// GetGlobalGGHomes is called by other functions within the cmd package to get the OGG Home path list
func GetGlobalGGHomes() []string {
	globalGGHomesMutex.RLock()
	defer globalGGHomesMutex.RUnlock()
	// Return a copy to prevent external modification
	copiedHomes := make([]string, len(globalGGHomes))
	copy(copiedHomes, globalGGHomes)
	return copiedHomes
}

