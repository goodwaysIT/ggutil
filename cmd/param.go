package cmd

import (
	"fmt"
	"os"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunParam handles the 'param' command to view parameter file content for specified processes.
func RunParam(c *cli.Context, processNames []string) error {
	fmt.Println("Executing 'param' command for processes:", processNames)
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	if len(processNames) == 0 {
		return cli.Exit("Error: No process names specified for 'param' command.", 1)
	}

	for _, home := range ggHomes {
		fmt.Printf("\n--- OGG Home: %s ---\n", home)
		for _, pName := range processNames {
			command := fmt.Sprintf("view param %s", pName)
			fmt.Printf("Executing: %s in %s\n", command, home)
			output, stderr, err := ogg.ExecuteGGSCICommand(home, command)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing '%s' in %s: %v\n", command, home, err)
				if output != "" {
					fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
				}
				if stderr != "" {
					fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderr)
				}
				continue // Continue with the next process or home
			}
			fmt.Printf("Output for 'view param %s':\n%s\n", pName, output)
			if stderr != "" {
				fmt.Printf("Stderr for 'view param %s':\n%s\n", pName, stderr)
			}
		}
	}
	return nil
}
