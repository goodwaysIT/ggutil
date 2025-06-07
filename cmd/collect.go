package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// RunCollect handles the 'collect' command to gather detailed information for specified processes.
func RunCollect(c *cli.Context, processNames []string) error {
	fmt.Println("Executing 'collect' command for processes:", processNames)
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	if len(processNames) == 0 {
		return cli.Exit("Error: No process names specified for 'collect' command.", 1)
	}

	for _, home := range ggHomes {
		fmt.Printf("\n--- OGG Home: %s ---\n", home)
		for _, pName := range processNames {
			fmt.Printf("Collecting details for process '%s':\n", pName)
		}
	}
	return nil
}
