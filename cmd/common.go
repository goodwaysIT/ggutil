package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
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

// executeGGSCICommand executes ggsci commands in the specified OGG Home.
// oggHome: Path to a single OGG Home.
// ggsciCommands: Command string(s) to execute in ggsci, multiple commands separated by newline.
func executeGGSCICommand(oggHome string, ggsciCommands string) (string, error) {
	ggsciPath := filepath.Join(oggHome, "ggsci.exe")

	// Check if ggsci.exe exists
	if _, err := os.Stat(ggsciPath); os.IsNotExist(err) {
		return "", fmt.Errorf("ggsci.exe not found at %s: %w", ggsciPath, err)
	}

	cmd := exec.Command(ggsciPath)

	// Set up standard input, output, and error
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdin pipe for ggsci: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start ggsci command at %s: %w", oggHome, err)
	}

	// Write ggsci commands
	// Ensure each command is followed by a newline for ggsci to execute it
	// Finally, add 'exit' command to ensure the ggsci process terminates
	commandsWithExit := ggsciCommands
	if !strings.HasSuffix(strings.TrimSpace(commandsWithExit), "\nexit") && !strings.Contains(strings.ToLower(commandsWithExit), "exit") {
		if !strings.HasSuffix(commandsWithExit, "\n") {
			commandsWithExit += "\n"
		}
		commandsWithExit += "exit\n"
	}

	_, err = io.WriteString(stdin, commandsWithExit)
	if err != nil {
		// Try to close stdin; if writing failed, this might help release resources or expose the issue
		stdin.Close() // Best effort
		// Try to wait for the command to get potential stderr output
		cmd.Wait() // Ignore error here as we are primarily concerned with the write failure
		return stdoutBuf.String(), fmt.Errorf("failed to write commands to ggsci stdin at %s (stderr: %s): %w", oggHome, stderrBuf.String(), err)
	}
	// Close stdin to signal end of command input
	if err := stdin.Close(); err != nil {
		// Even if closing stdin fails, try to wait for the command to finish
		cmd.Wait() // Best effort
		return stdoutBuf.String(), fmt.Errorf("failed to close ggsci stdin at %s (stderr: %s): %w", oggHome, stderrBuf.String(), err)
	}

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		// Command execution itself might return a non-0 exit code (e.g., 'info <proc>' if process doesn't exist).
		// stderrBuf might contain useful ggsci error messages in such cases.
		// Therefore, even if cmd.Wait() returns an error, we still return stdoutBuf and potentially stderrBuf content.
		// Append ggsci's stderr content to the error message for debugging.
		ggsciErrStr := strings.TrimSpace(stderrBuf.String())
		if ggsciErrStr != "" {
			return stdoutBuf.String(), fmt.Errorf("ggsci command execution failed at %s with exit error: %w. GGSCI Error: %s", oggHome, err, ggsciErrStr)
		}
		return stdoutBuf.String(), fmt.Errorf("ggsci command execution failed at %s with exit error: %w", oggHome, err)
	}

	// If stderr has content, it might be warnings or non-fatal errors. Append to result or log it.
	// For simplicity, only stdout is returned here. If needed, the function signature can be modified to return stderr.
	// if stderrBuf.Len() > 0 {
	// 	 fmt.Printf("Warning/Info from ggsci stderr at %s: %s\n", oggHome, stderrBuf.String())
	// }

	return stdoutBuf.String(), nil
}

// RunInfo handles the 'info' command
func RunInfo(c *cli.Context, processNames []string) error {
	fmt.Println("Executing 'info' command for processes:", processNames)
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	for _, home := range ggHomes {
		fmt.Printf("\n--- OGG Home: %s ---\n", home)
		for _, pName := range processNames {
			command := fmt.Sprintf("info %s\n", pName)
			output, err := executeGGSCICommand(home, command)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing 'info %s' in %s: %v\n", pName, home, err)
				// Depending on requirements, choose to continue with other processes/homes or return error directly.
				// continue // Continue to the next process or home
			}
			fmt.Printf("Output for 'info %s':\n%s\n", pName, output)
		}
	}
	return nil
}

// RunParam (Placeholder)
func RunParam(c *cli.Context, processNames []string) error {
	fmt.Println("Placeholder for 'param' command with processes:", processNames)
	fmt.Println("OGG Homes to be processed:", GetGlobalGGHomes())
	// TODO: Implement specific logic
	return nil
}

// RunConfig (Placeholder)
func RunConfig(c *cli.Context) error {
	fmt.Println("Placeholder for 'config' command.")
	fmt.Println("OGG Homes to be processed:", GetGlobalGGHomes())
	// TODO: Implement specific logic
	return nil
}

// RunBackup (Placeholder)
func RunBackup(c *cli.Context) error {
	fmt.Println("Placeholder for 'backup' command.")
	fmt.Println("OGG Homes to be processed:", GetGlobalGGHomes())
	// TODO: Implement specific logic
	return nil
}

// RunStats (Placeholder)
func RunStats(c *cli.Context, processName string, statArgs []string) error {
	fmt.Printf("Placeholder for 'stats' command for process '%s' with args: %v\n", processName, statArgs)
	fmt.Println("OGG Homes to be processed:", GetGlobalGGHomes())
	// TODO: Implement specific logic
	return nil
}

// RunCollect (Placeholder)
func RunCollect(c *cli.Context, processNames []string) error {
	fmt.Println("Placeholder for 'collect' command with processes:", processNames)
	fmt.Println("OGG Homes to be processed:", GetGlobalGGHomes())
	// TODO: Implement specific logic
	return nil
}
