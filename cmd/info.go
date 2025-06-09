package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunInfo handles the 'info' command for a single process name argument, querying all OGG Homes.
func RunInfo(c *cli.Context, processName string) error {
	ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Executing 'info' command for process '%s'\n", processName)

	gghomes := GetGlobalGGHomes()
	if len(gghomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}
	if processName == "" {
		return cli.Exit("Error: Process name is required for 'info' command.", 1)
	}

	// doInfo finds and prints info for the named process in one OGG Home.
	doInfo := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
		// Print debug info for the instance.
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Create new instance %d for %s: %v\n", id, home, gi)
		re := regexp.MustCompile(`[ ][ ]*`)
		line := strings.Split(re.ReplaceAllString(gi.GetInfoall(), " "), "\n")
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Infoall: %v\n", line)
		for _, v := range line {
			if strings.HasPrefix(v, "EXTRACT") || strings.HasPrefix(v, "REPLICAT") {
				field := strings.Fields(v)
				if field[2] == strings.ToUpper(processName) {
					er := ogg.NewGGER(home, v)
					ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Create new GGER from %s: %v\n", v, er)
					fmt.Fprintf(c.App.Writer, "\n%s\n", ogg.GGERInfo(&er, home))
				}
			}

		}
	}
	// Launch concurrent info fetch for each home.
	var wg sync.WaitGroup
	wg.Add(len(gghomes))
	for i, home := range gghomes {
		go doInfo(&wg, home, i)
	}
	wg.Wait()
	return nil
}
