package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/mholt/archiver/v3"
	"github.com/urfave/cli/v2"
)

// RunBackup handles the 'backup' command by listing important directories and files to be backed up.
func RunBackup(c *cli.Context) error {
	gghomes := GetGlobalGGHomes()
	destDir := filepath.Join("/tmp", "oggbackup_"+getHostname()+"_"+time.Now().Format("20060102_150405"))

	if len(gghomes) == 0 {
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "No OGG Home configured. Please specify using the -g parameter or GG_HOMES environment variable.\n")
		return nil
	}

	// doBackup copies files from one OGG Home.
	doBackup := func(wg *sync.WaitGroup, home string, id int) {
		defer wg.Done()
		gi := ogg.NewGGInst(home)
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Create new instance %d for %s: %v\n", id, home, gi)
		gi.SetMgr()
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Manager: %v\n", gi.Mgr)
		gi.SetER()
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "ER: %v\n", gi.Er)
		gi.SetFiles()
		// ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Files: %v\n", gi.jarFiles)
		for _, file := range gi.BackFileList() {
			// fmt.Println(file, destDir)
			ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "File: %s\n", file)
			copyFile(file, destDir)
		}
	}

	// Launch concurrent backup for each home.
	var wg sync.WaitGroup
	wg.Add(len(gghomes))
	for i, home := range gghomes {
		go doBackup(&wg, home, i)
	}
	wg.Wait()
	ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Backup to %s completed \n", destDir)
	TarGzFile(c, destDir)
	//remove folder
	if err := os.RemoveAll(destDir); err != nil {
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Error removing backup directory: %v\n", err)
	}
	return nil
}

// TarGzFile creates a tar.gz archive of the specified directory.
func TarGzFile(c *cli.Context, path string) {
	tarGz := archiver.NewTarGz()
	tarGz.SingleThreaded = false
	var source []string
	source = append(source, path)
	err := tarGz.Archive(source, path+".tar.gz")
	if err != nil {
		ogg.DebugPrint(c.Bool("debug"), c.App.Writer, "Error creating tar.gz: %v\n", err)
	}
	fmt.Printf("\nPlease refer to gz file %s\n", path+".tar.gz")
}
