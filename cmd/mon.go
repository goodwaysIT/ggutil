package cmd

import (
	"fmt"

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

		// Execute 'info all' command for each OGG Home
		output, err := executeGGSCICommand(homePath, "info all")
		if err != nil {
			fmt.Printf("Failed to execute 'info all' at %s: %v\n", homePath, err)
			// Continue to the next OGG Home even if one fails
			// If stdout contains partial information (handled in executeGGSCICommand), it will be printed here
			if output != "" {
			    fmt.Printf("Partial output:\n%s\n", output)
			}
			continue
		}

		// Print the output of the 'info all' command
		fmt.Println("GGSCI 'info all' command output:")
		fmt.Println(output)
	}

	fmt.Println("\n'mon' command execution completed.")
	return nil
}
