package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/mholt/archiver/v3"
	"github.com/urfave/cli/v2"
)

// RunCollect handles the 'collect' command to gather and zip all files related to a process for each OGG Home.
func RunCollect(c *cli.Context, processName string) error {
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}
	if processName == "" {
		return cli.Exit("Error: No process name specified for 'collect' command.", 1)
	}

	// doCollect collects and archives files for the specified process in a given home.
	doCollect := func(wg *sync.WaitGroup, home string, id int) {
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
					fileList, purgeList := er.CollectFileList(home)
					if len(fileList) == 0 {
						fmt.Fprintf(os.Stderr, "No files to collect for %s in %s\n", processName, home)
						return
					}
					// Compose archive name and path.
					zipName := fmt.Sprintf("oggcollect_%s_%s.zip", strings.ToLower(processName), time.Now().Format("20060102_150405"))
					zipFile := filepath.Join("/tmp", zipName)
					// Archive collected files into zip.
					err := archiver.Archive(fileList, zipFile)
					fmt.Println("\nPlease to refer to file -- " + zipFile + "\n")
					if err != nil {
						fmt.Println(err)
					}
					// Remove original files after archiving.
					for _, f := range purgeList {
						removeFile(f)
					}
					break
				}
			}
		}
	}

	// Launch concurrent collection for each home.
	var wg sync.WaitGroup
	wg.Add(len(ggHomes))
	for i, home := range ggHomes {
		go doCollect(&wg, home, i)
	}
	wg.Wait()
	return nil
}
