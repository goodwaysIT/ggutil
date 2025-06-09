package cmd

import (
	"fmt"
	"sync"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunMon is the handler for the 'mon' subcommand. It retrieves and prints monitoring information from all configured OGG instances.
func RunMon(c *cli.Context) error {
	gghomes := GetGlobalGGHomes() // Get the list of configured OGG Homes from common.go

	if len(gghomes) == 0 {
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "No OGG Home configured. Please specify using the -g parameter or GG_HOMES environment variable.\n")
		return nil
	}

	// doMon prints monitoring info for one OGG Home.
	doMon := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Create new instance %d for %s: %v\n", id, home, gi)
		fmt.Println(gi.Mon())
	}

	// Launch concurrent monitoring for each home.
	var wg sync.WaitGroup
	wg.Add(len(gghomes))
	for i, home := range gghomes {
		go doMon(&wg, home, i)
	}
	wg.Wait()
	return nil
}
