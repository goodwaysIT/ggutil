package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunStats handles the 'stats' command for a specific process.
func RunStats(c *cli.Context, processName string, statArgs []string) error {
	fmt.Printf("Executing 'stats' command for process '%s' with arguments: %v\n", processName, statArgs)
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	if processName == "" {
		return cli.Exit("Error: Process name is required for 'stats' command.", 1)
	}

	for _, home := range ggHomes {
		fmt.Printf("\n--- OGG Home: %s ---\n", home)
		var command string
		if len(statArgs) > 0 {
			command = fmt.Sprintf("stats %s, %s", processName, strings.Join(statArgs, ", "))
		} else {
			command = fmt.Sprintf("stats %s", processName)
		}
		fmt.Printf("Executing: %s\n", command)
		output, stderr, err := ogg.ExecuteGGSCICommand(home, command)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing '%s' in %s: %v\n", command, home, err)
				if output != "" {
					fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
				}
				if stderr != "" {
					fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderr)
				}
				continue
			}
			fmt.Printf("Output for '%s':\n%s\n", command, output)
			if stderr != "" {
				fmt.Printf("Stderr for '%s':\n%s\n", command, stderr)
			}
		}
	return nil
}
