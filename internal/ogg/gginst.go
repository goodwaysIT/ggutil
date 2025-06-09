package ogg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// GGInst represents a GoldenGate Home instance, encapsulating all processes, manager, and key files.
// Holds parsed process info, configuration, and file paths for one OGG Home.
type GGInst struct {
	soft       *GGSoft   // Pointer to OGG Home metadata
	Er         []GGER    // Extract and Replicat processes
	Mgr        GGMgr     // Manager process
	infoall    string    // Raw output of 'info all'
	globalFile string    // Path to GLOBALS file
	errLog     string    // Path to error log
	jarFiles   []string  // List of .jar files in dirprm
}

// GGInst methods for OGG instance management

func (gi *GGInst) setInfoall(home string) {
	gi.infoall, _ = execGGSCICmd(home, "info all")
}

func (gi *GGInst) SetMgr() {
	lines := strings.Split(gi.infoall, "\n")
	for _, v := range lines {
		if strings.HasPrefix(v, "MANAGER") {
			gi.Mgr = NewGGMgr(gi.soft.home, v)
			break
		}
	}
}

func (gi *GGInst) SetER() {
	var er []GGER
	var er1 GGER
	re := regexp.MustCompile(`[ ][ ]*`)
	line := strings.Split(re.ReplaceAllString(gi.infoall, " "), "\n")
	for _, v := range line {
		if strings.HasPrefix(v, "EXTRACT") || strings.HasPrefix(v, "REPLICAT") {
			er1 = NewGGER(gi.soft.home, v)
			er = append(er, er1)
		}
	}
	gi.Er = er
}

func (gi *GGInst) SetERConfig() {
	er := gi.Er
	var er1 *GGER
	var db2list []DB2List
	if gi.soft.ggType == "DB2" {
		db2list = getDB2List()
	}
	for i := 0; i < len(er); i++ {
		er1 = &er[i]
		er1.setConfig(gi.soft.home)
		if len(db2list) > 0 {
			er1.setDB2Alias(db2list)
		}
		er[i] = *er1
	}
	gi.Er = er
}

func (gi *GGInst) SetFiles() {
	gi.globalFile = filepath.Join(gi.soft.home, "GLOBALS")
	gi.errLog = filepath.Join(gi.soft.home, "ggserr.log")
	files, err := ioutil.ReadDir(filepath.Join(gi.soft.home, "dirprm"))
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jar") {
			gi.jarFiles = append(gi.jarFiles, filepath.Join(gi.soft.home, "dirprm", file.Name()))
		}
	}
}

// NewGGInst creates and initializes a GGInst instance
func NewGGInst(home string) *GGInst {
	var gi GGInst
	gi.soft = NewGGSoft(home)
	gi.setInfoall(home)
	// gi.setMgr()
	// gi.setER()
	// gi.setFiles()
	return &gi
}

func (gi *GGInst) Mon() string {
	res := "\n==== Home: " + gi.soft.home + ", " + gi.soft.oggFor + ", " + gi.soft.version + "\n" +
		gi.GetInfoall() +
		strings.Repeat("-", 80) + "\n"
	return res
}

func (gi *GGInst) GetInfoall() string {
	return gi.infoall
}

// RenderConfigTable returns a formatted config table for this OGG Home instance.
func (gi *GGInst) RenderConfigTable() string {
	res := "\n==== Home: " + gi.soft.home + ", " + gi.soft.oggFor + ", " + gi.soft.version + "\n\n"
	res += fmt.Sprintf("%-10s %-10s %-10s %-10s %-10s %-44s %-44s\n",
		"Program", "Status", "Group", "TabNo(prm)", "TabNo(rpt)", "Source", "Target")
	res += fmt.Sprintf("%-10s %-10s %-10s %-10s %-10s %-44s %-44s\n",
		strings.Repeat("-", 10), strings.Repeat("-", 10), strings.Repeat("-", 10),
		strings.Repeat("-", 10), strings.Repeat("-", 10),
		strings.Repeat("-", 44), strings.Repeat("-", 44))
	for _, v := range gi.Er {
		res += fmt.Sprintf("%-10s %-10s %-10s %-10d %-10d %-44s %-44s\n",
			v.program, v.status, v.group, v.prmTabCnt, v.rptTabCnt, v.source, v.target)
	}
	return res
}

func (gi *GGInst) BackFileList() []string {
	var fileList []string
	if isExist(gi.globalFile) {
		fileList = append(fileList, gi.globalFile)
	}
	if isExist(gi.errLog) {
		fileList = append(fileList, gi.errLog)
	}
	fileList = append(fileList, gi.jarFiles...)
	if isExist(gi.Mgr.rptFile) {
		fileList = append(fileList, gi.Mgr.rptFile)
	}
	if isExist(gi.Mgr.paramFile) {
		fileList = append(fileList, gi.Mgr.paramFile)
	}
	if len(gi.Mgr.obeyFiles) > 0 {
		fileList = append(fileList, gi.Mgr.obeyFiles...)
	}
	for _, v := range gi.Er {
		if isExist(v.rptFile) {
			fileList = append(fileList, v.rptFile)
		}
		if isExist(v.paramFile) {
			fileList = append(fileList, v.paramFile)
		}
		if len(v.obeyFiles) > 0 {
			fileList = append(fileList, v.obeyFiles...)
		}
		if v.propsFile != "" {
			fileList = append(fileList, v.propsFile)
		}
		if v.propertiesFile != "" {
			fileList = append(fileList, v.propertiesFile)
		}
	}
	return fileList
}
