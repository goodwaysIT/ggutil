package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bndr/gotabulate"
	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunTasks lists all SOURCEISTABLE tasks under all OGG instances, grouped by home, and prints them in formatted tables.
func RunTasks(c *cli.Context) error {
	// Retrieve all configured OGG Home directories.
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	// homeRows maps each home path to its list of SOURCEISTABLE task rows.
	var (
		homeRows = make(map[string][][]string)
		mtx      sync.Mutex // Protects homeRows in concurrent goroutines.
		wg       sync.WaitGroup
	)

	// Launch a goroutine for each OGG Home to query tasks concurrently.
	for _, home := range ggHomes {
		wg.Add(1)
		go func(home string) {
			defer wg.Done()
			// Run GGSCI command to list all tasks for this home.
			output, _, err := ogg.ExecuteGGSCICommand(home, "info *,tasks")
			if err != nil {
				fmt.Printf("[ERROR] Failed to run info *,tasks in %s: %v\n", home, err)
				return
			}
			// Find the start of each EXTRACT/REPLICAT block in the output.
			blocks := regexp.MustCompile(`(?m)^EXTRACT[ ].*|^REPLICAT[ ].*`).FindAllStringIndex(output, -1)
			var blockStarts []int
			for _, b := range blocks {
				blockStarts = append(blockStarts, b[0])
			}
			blockStarts = append(blockStarts, len(output)) // Add sentinel for last block.
			// Iterate over each process block.
			for i := 0; i < len(blockStarts)-1; i++ {
				block := output[blockStarts[i]:blockStarts[i+1]]
				// Only process blocks containing SOURCEISTABLE tasks.
				if !strings.Contains(block, "Task                 SOURCEISTABLE") {
					continue
				}
				lines := strings.Split(block, "\n")
				var prog, status, group, chkpt, task string
				// Parse each line in the process block.
				for _, line := range lines {
					f := strings.Fields(line)
					// Extract program type, group, and status from EXTRACT/REPLICAT header line.
					if len(f) >= 1 && (f[0] == "EXTRACT" || f[0] == "REPLICAT") {
						prog = f[0]
						if len(f) > 2 {
							group = f[1]
							status = f[len(f)-1]
						}
					}
					// Parse checkpoint line and its possible continuation.
					if strings.HasPrefix(line, "Log Read Checkpoint") {
						chkpt = strings.TrimSpace(strings.TrimPrefix(line, "Log Read Checkpoint"))
					} else if strings.HasPrefix(line, "                     ") && chkpt != "" {
						chkpt += " " + strings.TrimSpace(line)
					}
					// Detect SOURCEISTABLE task line.
					if strings.HasPrefix(line, "Task") && strings.Contains(line, "SOURCEISTABLE") {
						task = "SOURCEISTABLE"
					}
				}
				// If a SOURCEISTABLE task was found, add it to the result for this home.
				if task == "SOURCEISTABLE" {
					mtx.Lock()
					homeRows[home] = append(homeRows[home], []string{task, prog, status, group, chkpt})
					mtx.Unlock()
				}
			}
		}(home)
	}
	wg.Wait()

	// Print results for each home. If no tasks found, print an empty table.
	found := false
	for _, home := range ggHomes {
		rows := homeRows[home]
		fmt.Printf("\n==== OGG SOURCEISTABLE Tasks (%s) ====\n", home)
		if len(rows) == 0 {
			// No SOURCEISTABLE tasks in this home.
			fmt.Println()
			continue
		}
		// Use gotabulate to print the table for this home.
		tab := gotabulate.Create(rows)
		tab.SetHeaders([]string{"Task", "Program", "Status", "Group", "Checkpoint"})
		tab.SetAlign("left")
		fmt.Println(tab.Render("grid"))
		found = true
	}
	// If no tasks found in any home, print a summary message.
	if !found {
		fmt.Println("\nNo SOURCEISTABLE tasks found in any OGG instance.")
	}
	return nil
}
