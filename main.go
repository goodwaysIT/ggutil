package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goodwaysIT/ggutil/cmd" // Assumes cmd package is under ggutil module
	"github.com/goodwaysIT/ggutil/internal/ogg"
	"github.com/urfave/cli/v2"
)

// GlobalFlags stores the values of global flags
const Version = "1.0.0" // Application version
const RepoURL = "https://github.com/goodwaysIT/ggutil" // Open source repository URL

var GlobalFlags struct {
	GGHomes string
	Debug   bool
}

func main() {
	app := &cli.App{
		Name:  "ggutil",
		Usage: "Oracle GoldenGate multi-instance management tool\nOpen Source: " + RepoURL,
		// Define global flag -g, and attempt to get default value from GG_HOMES environment variable
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "gghomes",
				Aliases:     []string{"g"},
				Usage:       "Specify one or more OGG Home paths, comma-separated. If not specified, attempts to read from GG_HOMES environment variable.",
				EnvVars:     []string{"GG_HOMES"}, // Automatically read from GG_HOMES environment variable
				Destination: &GlobalFlags.GGHomes,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug output (show errors, warnings, exceptions)",
				Destination: &GlobalFlags.Debug,
			},
		},
		Before: func(c *cli.Context) error {
			// Process and validate OGG Home before any command execution
			if GlobalFlags.GGHomes == "" {
				// Re-check environment variable as urfave/cli's EnvVars behavior might need explicit handling
				envHomes := os.Getenv("GG_HOMES")
				if envHomes != "" {
					GlobalFlags.GGHomes = envHomes
				} else {
					// Set default OGG Home path if not specified
					defaultHome := "/acfsogg/oggo,/acfsogg/oggm,/acfsogg/oggp,/acfsogg/oggb" // Default value for OGG Home
					GlobalFlags.GGHomes = defaultHome
					if GlobalFlags.Debug {
						ogg.DebugPrint(GlobalFlags.Debug, nil, "Warning: OGG Home not specified, using default: %s\n", defaultHome)
					}
				}
			}
			if GlobalFlags.Debug {
				ogg.DebugPrint(GlobalFlags.Debug, nil, "Detected OGG Homes: %s\n", GlobalFlags.GGHomes)
			}
			parsedHomes := parseAndValidateGGHomes(GlobalFlags.GGHomes)
			if len(parsedHomes) == 0 {
				if GlobalFlags.Debug {
					ogg.DebugPrint(GlobalFlags.Debug, nil, "Error: Failed to parse any valid OGG Home paths.\n")
				}
				return cli.Exit("Error: Failed to parse any valid OGG Home paths.", 1)
			}
			cmd.SetGlobalGGHomes(parsedHomes)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show application version and open source repository",
				Action: func(c *cli.Context) error {
					fmt.Println("ggutil version:", Version)
					fmt.Println("Open Source Repository:", RepoURL)
					return nil
				},
			},
			{
				Name:      "tasks",
				Usage:     "List all OGG SOURCEISTABLE tasks under all homes.",
				ArgsUsage: "",
				Action: func(c *cli.Context) error {
					return cmd.RunTasks(c)
				},
			},
			{
				Name:  "mon",
				Usage: "Get version and path information for all OGG instances, print 'info all' results for each.",
				Action: func(c *cli.Context) error {
					return cmd.RunMon(c) // Call handler function in cmd package
				},
			},
			{
				Name:      "info",
				Usage:     "Get information for OGG processes (iterates over all configured OGG Homes).",
				ArgsUsage: "<process_name>", // Hint user that process name(s) are required
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.Exit("Error: 'info' command requires at least one process name argument.", 1)
					}
					processName := c.Args().First()
					return cmd.RunInfo(c, processName)
				},
			},
			{
				Name:      "param",
				Usage:     "Get parameter configuration for OGG processes (iterates over all configured OGG Homes).",
				ArgsUsage: "<process_name>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.Exit("Error: 'param' command requires at least one process name argument.", 1)
					}
					processName := c.Args().First()
					return cmd.RunParam(c, processName)
				},
			},
			{
				Name:  "config",
				Usage: "View process configuration details within OGG instances (iterates over all configured OGG Homes).",
				Action: func(c *cli.Context) error {
					return cmd.RunConfig(c)
				},
			},
			{
				Name:  "backup",
				Usage: "Backup configuration, log, report files, etc., for OGG instances (iterates over all configured OGG Homes).",
				Action: func(c *cli.Context) error {
					return cmd.RunBackup(c)
				},
			},
			{
				Name:      "stats",
				Usage:     "View statistics for a specific OGG process (total, daily, hourly) (iterates over all configured OGG Homes).",
				ArgsUsage: "<process_name>", // Process name is required, dimensions optional
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.Exit("Error: 'stats' command requires a process name argument.", 1)
					}
					processName := c.Args().First()
					// statArgs := c.Args().Slice()[1:] // Get arguments after process name as stat dimensions
					return cmd.RunStats(c, processName)
				},
			},
			{
				Name:      "collect",
				Usage:     "Collect information for a specific OGG process (info, infodetail, showch, status) (iterates over all configured OGG Homes).",
				ArgsUsage: "<process_name>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.Exit("Error: 'collect' command requires at least one process name argument.", 1)
					}
					processName := c.Args().First()
					return cmd.RunCollect(c, processName)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		if GlobalFlags.Debug {
			fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		}
		os.Exit(1)
	}
}

// parseAndValidateGGHomes parses and validates the OGG Home string.
// It converts a comma-separated (or semicolon-separated) string into a slice of strings,
// trimming whitespace from each path.
// Future enhancements could include checking if paths actually exist.
func parseAndValidateGGHomes(homes string) []string {
	if homes == "" {
		return []string{}
	}
	// Support both comma and semicolon as separators for flexibility
	replacedHomes := strings.ReplaceAll(homes, ";", ",")
	parts := strings.Split(replacedHomes, ",")
	validHomes := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			// Path validity check could be added here, e.g., os.Stat(trimmed)
			// For simplicity, only non-empty check for now
			validHomes = append(validHomes, trimmed)
		}
	}
	return validHomes
}

// Note: cmd.SetGlobalGGHomes, cmd.RunMon, etc., functions need to be defined in the cmd package.
// Usage of gotabulate will be implemented within respective cmd functions as needed.
