package ogg

import (
	"regexp"
	"strings"
)

// GGSoft holds basic OGG software info
// All comments and output in English

// GGSoft holds basic OGG software info for a given OGG Home.
// Includes home path, version string, OGG type, and display info.
type GGSoft struct {
	home    string // OGG Home directory
	v       string // Raw version string from GGSCI
	ggType  string // OGG type (e.g., Oracle, DB2, etc)
	oggFor  string // Display string for OGG flavor
	version string // Display version string
}

// setV sets the version string for GGSoft
func (gs *GGSoft) setV() {
	gs.v, _ = execShell(gs.home + "/ggsci -v")
}

// setType sets the ggType field for GGSoft
func (gs *GGSoft) setType() {
	re := regexp.MustCompile(`[ ][ ]*`)
	line := strings.Split(re.ReplaceAllString(gs.v, " "), "\n")
	for _, v := range line {
		if strings.Contains(v, " for ") {
			s := strings.Split(v, " for ")
			gs.ggType = s[1]
			break
		}
		// for Oracle 12c/19c
		if strings.Contains(v, "Oracle 1") {
			oracleVersion := strings.Split(strings.Split(v, " on ")[0], ",")[3]
			gs.ggType = gs.ggType + oracleVersion
		}
	}
}

// setFor sets the oggFor and version fields for GGSoft
func (gs *GGSoft) setFor() {
	var oggFor, version string
	line := strings.Split(gs.v, "\n")
	for _, v := range line {
		if strings.HasPrefix(v, "Oracle GoldenGate") && strings.Contains(v, " for ") {
			oggFor = strings.Replace(v, " Command Interpreter", "", 1)
			oggFor = strings.Replace(oggFor, "Oracle GoldenGate", "OGG", 1)
		}
		if strings.HasPrefix(v, "Version") && !strings.Contains(v, "Build") {
			version = v
		}
		if oggFor != "" && version != "" {
			break
		}
	}
	gs.oggFor, gs.version = oggFor, version
}

// NewGGSoft creates and initializes a GGSoft instance
func NewGGSoft(home string) *GGSoft {
	var gs GGSoft
	gs.home = home
	gs.setV()
	gs.setType()
	gs.setFor()
	return &gs
}
