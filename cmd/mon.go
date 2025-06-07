package cmd

import (
	"fmt"
	"os"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunMon is the handler for the 'mon' subcommand.
// It retrieves information from all configured OGG instances.
func RunMon(c *cli.Context) error {
	homes := GetGlobalGGHomes() // Get the list of configured OGG Homes from common.go

	if len(homes) == 0 {
		fmt.Println("No OGG Home configured. Please specify using the -g parameter or GG_HOMES environment variable.")
		return nil // Or return an error cli.Exit("No OGG Home configured", 1)
	}

	fmt.Println("Starting 'mon' command to monitor all configured OGG instances...")

	for _, homePath := range homes {
		fmt.Printf("\n--- OGG Home: %s ---\n", homePath)

		// Get OGG version and path information
		// This might involve running a command like 'info all' or checking specific files
		// For simplicity, let's assume 'info all' gives enough details for now.
		fmt.Printf("Executing 'info all' in %s\n", homePath)
		output, stderr, err := ogg.ExecuteGGSCICommand(homePath, "info all")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing 'info all' in %s: %v\n", homePath, err)
			if output != "" {
				fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
			}
			if stderr != "" {
				fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderr)
			}
			continue // Try next home if error occurs
		}
		fmt.Printf("OGG Home: %s, Version and Path Info:\n%s\n", homePath, output)
		if stderr != "" {
			fmt.Printf("Stderr for 'info all' in %s:\n%s\n", homePath, stderr)
		}
		// Print the output of the 'info all' command
		fmt.Println("GGSCI 'info all' command output:")
		fmt.Println(output)
	}

	fmt.Println("\n'mon' command execution completed.")
	return nil
}
