package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunParam handles the 'param' command to view parameter file content for specified processes in all OGG Homes.
func RunParam(c *cli.Context, processName string) error {
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	if processName == "" {
		return cli.Exit("Error: No process names specified for 'param' command.", 1)
	}

	re := regexp.MustCompile(`[ ][ ]*`)
	// doParam finds and prints parameter file content for the named process in one OGG Home.
	doParam := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
		line := strings.Split(re.ReplaceAllString(gi.GetInfoall(), " "), "\n")
		for _, v := range line {
			if strings.HasPrefix(v, "EXTRACT") || strings.HasPrefix(v, "REPLICAT") {
				field := strings.Fields(v)
				if field[2] == strings.ToUpper(processName) {
					er := ogg.NewGGER(home, v)
					fmt.Printf("\n==== OGG Process [ %s ] Under Home: [ %s ] ====\n\n", field[2], home)
					paramFile, paramContent := ogg.GGERParam(&er)
					fmt.Printf("Param file [ %s ] content for '%s':\n\n%s\n", paramFile, field[2], paramContent)
				}
			}
		}
	}
	// Launch concurrent parameter viewing for each home.
	var wg sync.WaitGroup
	wg.Add(len(ggHomes))
	for i, home := range ggHomes {
		go doParam(&wg, home, i)
	}
	wg.Wait()
	return nil
}
