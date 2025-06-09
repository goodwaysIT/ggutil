package ogg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExecuteGGSCICommand executes ggsci commands in the specified OGG Home.
// oggHome: Path to a single OGG Home.
// ggsciCommands: Command string(s) to execute in ggsci, multiple commands separated by newline.
// It returns stdout, stderr, and an error if the command fails to execute or returns a non-zero exit code.
func ExecuteGGSCICommand(oggHome string, ggsciCommands string) (string, string, error) {
	ggsciPath := filepath.Join(oggHome, "ggsci")

	// Check if ggsci.exe exists
	if _, err := os.Stat(ggsciPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("ggsci not found at %s: %w", ggsciPath, err)
	}

	cmd := exec.Command(ggsciPath)

	// Set up standard input, output, and error
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", "", fmt.Errorf("failed to get stdin pipe for ggsci: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", stderrBuf.String(), fmt.Errorf("failed to start ggsci command at %s: %w", oggHome, err)
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
		stdin.Close() // Best effort
		cmd.Wait()    // Ignore error
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("failed to write commands to ggsci stdin at %s: %w", oggHome, err)
	}
	if err := stdin.Close(); err != nil {
		cmd.Wait() // Best effort
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("failed to close ggsci stdin at %s: %w", oggHome, err)
	}

	if err := cmd.Wait(); err != nil {
		return stdoutBuf.String(), stderrBuf.String(), fmt.Errorf("ggsci command execution failed at %s with exit error: %w", oggHome, err)
	}
	return stdoutBuf.String(), stderrBuf.String(), nil
}
