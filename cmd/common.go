package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	// "github.com/urfave/cli/v2" // No longer needed here if all Run... functions are moved
)

// globalGGHomes stores the list of OGG Home paths set by the main package.
var globalGGHomes []string
var globalGGHomesMutex sync.RWMutex

// SetGlobalGGHomes sets the parsed OGG Home path list (thread-safe, called by main).
func SetGlobalGGHomes(homes []string) {
	globalGGHomesMutex.Lock()
	defer globalGGHomesMutex.Unlock()
	globalGGHomes = homes
}

// GetGlobalGGHomes retrieves the OGG Home path list (thread-safe, returns a copy).
func GetGlobalGGHomes() []string {
	globalGGHomesMutex.RLock()
	defer globalGGHomesMutex.RUnlock()
	// Return a copy to prevent external modification.
	copiedHomes := make([]string, len(globalGGHomes))
	copy(copiedHomes, globalGGHomes)
	return copiedHomes
}

// getHostname returns the current system hostname.
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	return hostname
}

// createDir creates the specified directory if it does not exist.
func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0777)
		// log.Info("Create directory " + dir)
	}
}

// copyFile copies a file to a destination directory, preserving subdirectory structure.
func copyFile(file, destDir string) error {
	in, err := os.Open(file)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return err
	}
	defer in.Close()

	// Construct the actual destination directory, preserving the original structure.
	actualDir := strings.Replace(filepath.Dir(file), "/", destDir+"/", 1)
	createDir(actualDir)
	destFile := filepath.Join(actualDir, filepath.Base(in.Name()))
	// fmt.Println(destFile)
	out, err := os.Create(destFile)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return err
	}
	defer out.Close()

	// Copy file contents.
	if _, err := io.Copy(out, in); err != nil {
		fmt.Printf("err %v\n", err)
		return err
	}
	return out.Close()
}

// removeFile deletes a file and prints errors if any.
func removeFile(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
	}
}

// getBasedir returns the absolute path of the application's base directory.
func getBasedir() string {
	programPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	// programName := path.Base(os.Args[0])
	return programPath
}
