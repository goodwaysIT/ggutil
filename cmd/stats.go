package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunStats handles the 'stats' command for a specific process in all OGG Homes.
func RunStats(c *cli.Context, processName string) error {
	ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Executing 'stats' command for process '%s'\n", processName)
	gghomes := GetGlobalGGHomes()
	if len(gghomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	if processName == "" {
		return cli.Exit("Error: Process name is required for 'stats' command.", 1)
	}

	// doStats finds and prints stats for the named process in one OGG Home.
	doStats := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
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
					fmt.Printf("\n==== OGG Process [ %s ] Under Home: [ %s ] ====\n\n", field[2], home)
					fmt.Println(strings.Repeat("=", 40) + "[total stats]" + strings.Repeat("=", 40))
					fmt.Println(ogg.GGERStats(&er, home, "total"))
					fmt.Println(strings.Repeat("=", 40) + "[daily stats]" + strings.Repeat("=", 40))
					fmt.Println(ogg.GGERStats(&er, home, "daily"))
					fmt.Println(strings.Repeat("=", 40) + "[hourly stats/sec]" + strings.Repeat("=", 40))
					fmt.Println(ogg.GGERStats(&er, home, "hourly, reportrate sec"))
				}
			}

		}
	}

	// Launch concurrent stats viewing for each home.
	var wg sync.WaitGroup
	wg.Add(len(gghomes))
	for i, home := range gghomes {
		go doStats(&wg, home, i)
	}
	wg.Wait()
	return nil
}
