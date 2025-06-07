package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// RunBackup handles the 'backup' command by listing important directories and files to be backed up.
func RunBackup(c *cli.Context) error {
	fmt.Println("Executing 'backup' command (listing items to back up, no actual backup performed).")
	ggHomes := GetGlobalGGHomes()
	if len(ggHomes) == 0 {
		return cli.Exit("Error: OGG Home list is empty. Please check configuration.", 1)
	}

	for _, home := range ggHomes {
		ggsciCommands := "INFO ALL\nSEND MANAGER GETPARAMS"
		fmt.Printf("Attempting to get information for backup from OGG Home: %s\n", home)
		output, stderr, err := ogg.ExecuteGGSCICommand(home, ggsciCommands)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing backup-related commands in %s: %v\n", home, err)
			if output != "" {
				fmt.Fprintf(os.Stderr, "Stdout:\n%s\n", output)
			}
			if stderr != "" {
				fmt.Fprintf(os.Stderr, "Stderr:\n%s\n", stderr)
			}
			continue
		}
		fmt.Printf("--- Backup Information for %s ---\n%s\n", home, output)
		if stderr != "" {
			fmt.Printf("Stderr for backup commands in %s:\n%s\n", home, stderr)
		}
		fmt.Println("Recommended items to back up for this OGG Home:")
		fmt.Printf("  - Parameter files: %s\n", filepath.Join(home, "dirprm"))
		fmt.Printf("  - Report files: %s\n", filepath.Join(home, "dirrpt"))
		fmt.Printf("  - Checkpoint files: %s\n", filepath.Join(home, "dirchk"))
		fmt.Printf("  - Definition files: %s\n", filepath.Join(home, "dirdef"))
		fmt.Printf("  - Main error log: %s\n", filepath.Join(home, "ggserr.log"))
		fmt.Println("  - Wallet files (if applicable, e.g., Oracle SSL/TLS):", filepath.Join(home, "dirwlt"))
		fmt.Println("  - Other critical directories like 'BR' (Backup/Recovery), 'cfg', etc. as per your setup.")
		// Note: Actual backup implementation (e.g., zipping files) is not done here.
	}
	return nil
}
