package cmd

import (
	"fmt"
	"sync"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunConfig handles the 'config' command to view parameter files for all major processes in all OGG Homes.
func RunConfig(c *cli.Context) error {
	gghomes := GetGlobalGGHomes()

	if len(gghomes) == 0 {
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "No OGG Home configured. Please specify using the -g parameter or GG_HOMES environment variable.\n")
		return nil
	}

	// doConfig fetches and prints config table for one OGG Home.
	doConfig := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Create new instance %d for %s: %v\n", id, home, gi)
		gi.SetER()
		gi.SetERConfig()
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "ER Config: %v\n", gi.Er)
		fmt.Println(gi.RenderConfigTable())
	}

	// Launch concurrent config fetch for each home.
	var wg sync.WaitGroup
	wg.Add(len(gghomes))
	for i, home := range gghomes {
		go doConfig(&wg, home, i)
	}
	wg.Wait()
	return nil
}
