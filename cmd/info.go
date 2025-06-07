package cmd

import (
	"fmt"
	"os"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

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
			fmt.Printf("Executing: %s in %s", command, home)
			output, stderr, err := ogg.ExecuteGGSCICommand(home, command)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing 'info %s' in %s: %v\n", pName, home, err)
				if output != "" {
					fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
				}
				if stderr != "" {
					fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderr)
				}
				continue // Continue with the next process or home
			}
			fmt.Printf("Output for 'info %s' in %s:\n%s\n", pName, home, output)
			if stderr != "" {
				fmt.Printf("Stderr for 'info %s' in %s:\n%s\n", pName, home, stderr)
			}
		}
	}
	return nil
}
