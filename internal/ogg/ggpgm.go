// Package ogg provides core types and utilities for Oracle GoldenGate CLI integration.
package ogg

import (
	"path/filepath"
	"strings"
)

// GGPgm represents a generic OGG program/process with its status and related files.
type GGPgm struct {
	program   string   // Program name (EXTRACT, REPLICAT, etc)
	status    string   // Current status
	rptFile   string   // Path to report file
	paramFile string   // Path to parameter file
	obeyFiles []string // List of obey files referenced in param file
}

// GGMgr represents the OGG Manager process, embedding generic program fields.
type GGMgr struct {
	GGPgm // Embedded generic program fields
	// Manager specific fields can be added here
}

// NewGGMgr creates and initializes a GGMgr (Manager) instance from GGSCI info output.
func NewGGMgr(home, info string) GGMgr {
	var mgr GGMgr
	field := strings.Fields(info)
	if len(field) == 2 {
		mgr.program = field[0]
		mgr.status = field[1]
	}
	// Set report and parameter file paths for Manager
	mgr.rptFile = filepath.Join(home, "/dirrpt/MGR.rpt")
	mgr.paramFile = filepath.Join(home, "/dirprm/mgr.prm")
	return mgr
}
