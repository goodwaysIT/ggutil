package ogg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// execShell executes a shell command and returns its output as string
func execShell(cmd string) (string, error) {
	execCmd := exec.Command("/bin/bash", "-c", cmd)
	res, err := execCmd.Output()
	if err != nil {
		// panic(err)
		return string(res), errors.New("execCmd.Output error")
	} else {
		return string(res), nil
	}
}

// execGGSCICmd runs a GGSCI command in the given OGG home and returns its output.
func execGGSCICmd(home, command string) (string, error) {
	return execShell(`echo "` + command + `" | ` + home + `/ggsci -s | grep -v ^GGSCI `)
}

// readLines reads a file and returns its lines as a string slice.
func readLines(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

// isExist checks if a file exists at the given path.
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// writeContent writes the given content to a file at the specified path.
func writeContent(content, file string) error {
	return os.WriteFile(file, []byte(content), 0644)
}

// UniqueStrings helper function to remove duplicates from a slice of strings.
func UniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// DebugPrint prints debug messages if debug flag is true.
// If writer is nil, it defaults to os.Stdout.
// Usage: ogg.DebugPrint(GlobalFlags.Debug, c.App.Writer, "msg: %v", err)
func DebugPrint(debug bool, writer io.Writer, format string, a ...interface{}) {
	if debug {
		if writer == nil {
			writer = os.Stdout
		}
		fmt.Fprintf(writer, format, a...)
	}
}
