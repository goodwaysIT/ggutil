package cmd

import (
	"fmt"
	"os"

	"github.com/goodwaysIT/ggutil/internal/ogg" // Added for ogg.ExecuteGGSCICommand and ogg.ParseProcessNamesFromInfoAll
	"github.com/urfave/cli/v2"
)

// RunConfig handles the 'config' command to view parameter files for all major processes.
func RunConfig(c *cli.Context) error {
	fmt.Println("Executing 'config' command to view parameter files of major processes.")
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	for _, home := range ggHomes {
		fmt.Printf("\n--- OGG Home: %s ---\n", home)
		fmt.Println("Fetching all process information...")
		infoAllOutput, infoAllStderr, err := ogg.ExecuteGGSCICommand(home, "info all")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing 'info all' in %s: %v\n", home, err)
			if infoAllOutput != "" {
				fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", infoAllOutput)
			}
			if infoAllStderr != "" {
				fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", infoAllStderr)
			}
			continue
		}

		processNames := ogg.ParseProcessNamesFromInfoAll(infoAllOutput)
		if len(processNames) == 0 {
			fmt.Println("No processes found or parsed from 'info all' output.")
			continue
		}
		fmt.Printf("Found processes: %v\n", processNames)

		for _, pName := range processNames {
			command := fmt.Sprintf("view param %s", pName)
			fmt.Printf("Executing: %s\n", command)
			output, stderrOutput, err := ogg.ExecuteGGSCICommand(home, command)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing '%s' in %s: %v\n", command, home, err)
				if output != "" {
					fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
				}
				if stderrOutput != "" {
					fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderrOutput)
				}
				continue // Continue with the next process or home
			}
			fmt.Printf("Output for 'view param %s':\n%s\n", pName, output)
			if stderrOutput != "" {
				fmt.Printf("Stderr for 'view param %s':\n%s\n", pName, stderrOutput)
			}
		}
	}
	return nil
}
