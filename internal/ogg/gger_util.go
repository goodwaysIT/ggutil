package ogg

import (
	"fmt"
	"strings"

	"github.com/bndr/gotabulate"
)

// Stats holds parsed statistics for a table from GGSCI stats output.
// Each field corresponds to a column in the GoldenGate statistics output.
type Stats struct {
	table     string   // Table name
	insert    string   // Number of inserts
	update    string   // Number of updates
	before    string   // Number of befores
	delete    string   // Number of deletes
	upsert    string   // Number of upserts
	discard   string   // Number of discards
	operation string   // Number of operations
}

// GGERStats parses and formats the output of getStats for a GGER instance.
// It extracts table-level statistics from the command output and formats them as a table.
func GGERStats(er *GGER, home, interval string) string {
	// Retrieve raw stats output from GGSCI
	statsRes := er.getStats(home, interval)
	lines := strings.Split(statsRes, "\n")
	var stats []Stats
	var s Stats
	since := ""
	// Parse each line for table and operation stats
	for _, line := range lines {
		// Identify table name
		if strings.HasPrefix(line, "Extracting from") || strings.HasPrefix(line, "Replicating from") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				s.table = fields[2]
			}
		}
		// Capture the 'statistics since' line
		if since == "" && strings.Contains(line, " statistics since") {
			since = line
		}
		// Parse operation totals
		if strings.Contains(line, "Total") && !strings.Contains(line, "statistics") {
			l := strings.Fields(line)
			if len(l) < 3 {
				continue
			}
			res := l[2]
			switch {
			case strings.Contains(l[1], "inserts"):
				s.insert = res
			case strings.Contains(l[1], "updates"):
				s.update = res
			case strings.Contains(l[1], "befores"):
				s.before = res
			case strings.Contains(l[1], "deletes"):
				s.delete = res
			case strings.Contains(l[1], "upserts"):
				s.upsert = res
			case strings.Contains(l[1], "discards"):
				s.discard = res
			case strings.Contains(l[1], "operations"):
				s.operation = res
				stats = append(stats, s) // Save completed stats row
				s = Stats{} // Reset for next table
			}
		}
	}

	var sb strings.Builder
	// Print the statistics since line if present
	if since != "" {
		sb.WriteString("\n" + since + "\n")
	}
	// Prepare and render table output
	title := []string{"Table Name", "Insert", "Updates", "Befores", "Deletes", "Upserts", "Discards", "Operations"}
	var resRows [][]string
	for _, v := range stats {
		row := []string{v.table, v.insert, v.update, v.before, v.delete, v.upsert, v.discard, v.operation}
		resRows = append(resRows, row)
	}
	tab := gotabulate.Create(resRows)
	tab.SetHeaders(title)
	tab.SetAlign("left")
	sb.WriteString(tab.Render("grid"))
	return sb.String()
}

// GGERParam reads the prm file for the given GGER and returns its contents as a string.
// Returns the file path and its contents, or an error string if reading fails.
func GGERParam(er *GGER) (string, string) {
	if er.paramFile == "" {
		return "", "[ERROR] No param file path found for this process."
	}
	content, err := readLines(er.paramFile)
	if err != nil {
		return er.paramFile, fmt.Sprintf("[ERROR] Failed to read param file %s: %v", er.paramFile, err)
	}
	return er.paramFile, strings.Join(content, "\n")
}

// GGERInfo returns a detailed string for a GGER instance, including info, detail, showch, and file attributes.
// Formats the main process attributes as a table and lists obey files if present.
func GGERInfo(er *GGER, home string) string {
	fmt.Printf("\n==== OGG Process [ %s ] Under Home: [ %s ] ====\n", er.group, home)
	var sb strings.Builder
	// Basic info as table
	data := [][]interface{}{
		{"Group", er.group},
		{"Program", er.program},
		{"Status", er.status},
		{"Lag", er.lag},
		{"ChkptTime", er.chkptTime},
		{"Report File", er.rptFile},
		{"Param File", er.paramFile},
	}
	// Render the main info table
	tab := gotabulate.Create(data)
	tab.SetHeaders([]string{"Field", "Value"})
	tab.SetWrapStrings(false)
	tab.SetAlign("left")
	sb.WriteString(tab.Render("grid"))

	if len(er.obeyFiles) > 0 {
		sb.WriteString("obeyFiles:\n")
		for _, ob := range er.obeyFiles {
			sb.WriteString("  - " + ob + "\n")
		}
	}
	if er.propsFile != "" {
		sb.WriteString(fmt.Sprintf("Props File: %s\n", er.propsFile))
	}
	if er.propertiesFile != "" {
		sb.WriteString(fmt.Sprintf("Properties File: %s\n", er.propertiesFile))
	}
	// Info sections
	sb.WriteString("\n" + strings.Repeat("=", 40) + "[info detail]" + strings.Repeat("=", 40) + "\n")
	sb.WriteString(er.getInfoDetail(home))
	sb.WriteString("\n" + strings.Repeat("=", 40) + "[info showch]" + strings.Repeat("=", 40) + "\n")
	sb.WriteString(er.getInfoShowch(home))
	return sb.String()
}
