package ogg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// GGER represents an OGG Extract or Replicat process, including config and runtime attributes.
// Embedded GGPgm provides common program fields. Additional fields track files, stats, and process state.
type GGER struct {
	GGPgm                 // Embedded struct for program name/status
	group          string // Process group name
	lag            string // Lag value
	chkptTime      string // Checkpoint time
	propsFile      string // Path to .props file (if any)
	propertiesFile string // Path to .properties file (if any)
	prmTabCnt      int    // Table count from prm file
	rptTabCnt      int    // Table count from rpt file
	source         string // Source DB/table
	target         string // Target DB/table
}

func (er *GGER) getInfo(home string) string {
	var res string
	res, _ = execGGSCICmd(home, "info "+er.program+" "+er.group)
	return res
}

func (er *GGER) getInfoDetail(home string) string {
	var res string
	res, _ = execGGSCICmd(home, "info "+er.program+" "+er.group+", detail")
	return res
}

func (er *GGER) getInfoShowch(home string) string {
	var res string
	res, _ = execGGSCICmd(home, "info "+er.program+" "+er.group+", showch")
	return res
}

func NewGGER(home, info string) GGER {
	var er1 GGER
	field := strings.Fields(info)
	if len(field) == 5 {
		er1.program = field[0]
		er1.status = field[1]
		er1.group = field[2]
		er1.lag = field[3]
		er1.chkptTime = field[4]
		er1.rptFile = filepath.Join(home, `/dirrpt/`+field[2]+`.rpt`)
		er1.paramFile = filepath.Join(home, `/dirprm/`+strings.ToLower(field[2])+`.prm`)

		//other file
		lines, err := readLines(er1.paramFile)
		if err != nil {
			DebugPrint(IsDebugMode, nil, "%v\n", err)
		}
		var obFile string
		var propsFile, propertiesFile string

		for _, v := range lines {
			//obey file
			if strings.HasPrefix(v, "obey") || strings.HasPrefix(v, "OBEY") {
				f := strings.Split(v, " ")
				if strings.HasPrefix(f[1], "./") {
					obFile = filepath.Join(home, strings.Replace(f[1], "./", "/", -1))
				} else {
					obFile = f[1]
				}
				er1.obeyFiles = append(er1.obeyFiles, obFile)
			}
			//props file
			if (!strings.HasPrefix(v, "--")) && strings.Contains(v, "property") {
				f := strings.Split(v, "=")
				if strings.HasPrefix(f[1], "./") {
					propsFile = filepath.Join(home, strings.Replace(f[1], "./", "/", -1))
				} else {
					propsFile = filepath.Join(home, f[1])
				}
			}
		}
		//properties file
		if propsFile != "" {
			lines, err := readLines(propsFile)
			if err != nil {
				DebugPrint(IsDebugMode, nil, "%v\n", err)
			}
			for _, v := range lines {
				if (!strings.HasPrefix(v, "#")) && strings.HasSuffix(v, "properties") {
					f := strings.Split(v, "=")
					if len(f) > 1 {
						propertiesFile = filepath.Join(home, "dirprm", f[1])
					}
				}
			}
		}
		er1.propsFile, er1.propertiesFile = propsFile, propertiesFile
	}

	return er1
}

func (er *GGER) getStats(home, interval string) string {
	res, _ := execGGSCICmd(home, "stats "+er.program+" "+er.group+", "+interval)
	// fmt.Println(res)
	return res
}

func (er *GGER) getSendStatus(home string) string {
	res, _ := execGGSCICmd(home, "send "+er.program+" "+er.group+" status")
	// fmt.Println(res)
	return res
}

func (er *GGER) setConfig(home string) {
	var source, target string
	var t1, t2 int
	//prm file

	if er.program == "EXTRACT" {
		lines, err := readLines(er.paramFile)
		if err != nil {
			DebugPrint(IsDebugMode, nil, "%v\n", err)
		}
		for _, v := range lines {
			if strings.HasPrefix(v, "--") {
				continue
			}

			// for Oracle DB
			if strings.HasPrefix(v, "userid") || strings.HasPrefix(v, "USERID") {
				field := strings.Fields(v)
				s1 := strings.TrimRight(field[1], ",")
				s := strings.Split(s1, "@")
				if len(s) > 1 {
					source = s[1]
					if !strings.Contains(s[1], "/") {
						source = source + "(" + getOracleInfo(s[1]) + ")"
					}
				} else {
					source = s1
				}

				break
			}
			//for Non-Oracle DB
			if strings.HasPrefix(v, "sourcedb") || strings.HasPrefix(v, "SOURCEDB") {
				field := strings.Split(v, " ")
				source = strings.TrimRight(field[1], ",")
				break
			}

			//local pump
			if strings.HasPrefix(v, "exttrail") || strings.HasPrefix(v, "EXTTRAIL") {
				field := strings.Fields(v)
				target = field[1]
				break
			}

			//remote pump
			if strings.HasPrefix(v, "rmthost") || strings.HasPrefix(v, "RMTHOST") {
				field := strings.Fields(v)
				target = strings.TrimRight(field[1], ",")
				continue
			}
			if strings.HasPrefix(v, "rmtfile") || strings.HasPrefix(v, "RMTFILE") {
				field := strings.Fields(v)
				target = target + ":" + field[1]
				continue
			}
		}

		if target == "" {
			//SEQ & RBA
			var trail, seqno, rba string
			detail := er.getInfoDetail(home)
			lines := strings.Split(detail, "\n")
			for _, v := range lines {
				if strings.Contains(v, "EXTTRAIL") {
					field := strings.Fields(v)
					trail = field[0]
					s, _ := strconv.Atoi(field[1])
					seqno = fmt.Sprintf("%09d", s)
					rba = field[2]
					break
				}
			}
			target = trail + seqno + "(" + rba + ")"
		}

		//For PUMP Extract
		if source == "" {
			var seqno, rba string
			info := er.getInfo(home)
			lines := strings.Split(info, "\n")
			for _, v := range lines {
				if strings.Contains(v, "Log Read Checkpoint  File") {
					field := strings.Fields(v)
					seqno = field[4]
				}
				if strings.Contains(v, "RBA ") {
					field := strings.Fields(v)
					rba = field[3]
				}
			}
			source = seqno + "(" + rba + ")"
			//Change Program from EXTRACT to EXTRACT*p
			er.program = er.program + "*p"
		}
	}
	if er.program == "REPLICAT" {
		var seqno, rba string
		detail := er.getInfoDetail(home)
		lines := strings.Split(detail, "\n")
		for _, v := range lines {
			if strings.Contains(v, "File") {
				field := strings.Fields(v)
				seqno = field[4]
				continue
			}
			if strings.Contains(v, "RBA") {
				field := strings.Fields(v)
				rba = field[3]
				continue
			}
		}
		source = seqno + "(" + rba + ")"

		lines, err := readLines(er.paramFile)
		if err != nil {
			DebugPrint(IsDebugMode, nil, "%v\n", err)
		}
		for _, v := range lines {
			if strings.HasPrefix(v, "--") {
				continue
			}
			//for  DB
			if strings.HasPrefix(v, "targetdb") || strings.HasPrefix(v, "TARGETDB") {
				field := strings.Fields(v)
				target = strings.TrimRight(field[1], ",")
			}
			// for Oracle DB
			if strings.HasPrefix(v, "userid") || strings.HasPrefix(v, "USERID") {
				field := strings.Fields(v)
				s1 := strings.TrimRight(field[1], ",")
				s := strings.Split(s1, "@")
				if len(s) > 1 {
					target = s[1]
					if !strings.Contains(s[1], "/") {
						target = target + "(" + getOracleInfo(s[1]) + ")"
					}
				} else {
					target = s1
				}
				continue
			}
			// for Bigdata
			if strings.HasPrefix(v, "targetdb libfile") || strings.HasPrefix(v, "TARGETDB LIBFILE") {
				if strings.HasPrefix(er.group, "ES") {
					target = "ElasticSearch"
				} else {
					target = "Kafka"
				}
				continue
			}
		}

	}

	//prm file
	prmLine, err := readLines(er.paramFile)
	if err != nil {
		DebugPrint(IsDebugMode, nil, "%v\n", err)
	}
	for _, v := range prmLine {
		//table count
		if strings.HasPrefix(v, "map") || strings.HasPrefix(v, "MAP") ||
			strings.HasPrefix(v, "table") || strings.HasPrefix(v, "TABLE") {
			t1 = t1 + 1
			continue
		}
	}
	// obey file
	if len(er.obeyFiles) > 0 {
		for _, v := range er.obeyFiles {
			if isExist(v) {
				lines, err := readLines(v)
				if err != nil {
					DebugPrint(IsDebugMode, nil, "%v\n", err)
				}
				for _, v := range lines {
					//table count
					if strings.HasPrefix(v, "map") || strings.HasPrefix(v, "MAP") ||
						strings.HasPrefix(v, "table") || strings.HasPrefix(v, "TABLE") {
						t1 = t1 + 1
						continue
					}
				}
			}
		}
	}
	//rpt file
	rptLine, err := readLines(er.rptFile)
	if err != nil {
		DebugPrint(IsDebugMode, nil, "%v\n", err)
	}
	for _, v := range rptLine {
		//table count
		if strings.HasPrefix(v, "map") || strings.HasPrefix(v, "MAP") ||
			strings.HasPrefix(v, "table") || strings.HasPrefix(v, "TABLE") {
			t2 = t2 + 1
			continue
		}
		if strings.Contains(v, "Run Time Messages") {
			break
		}
	}

	er.prmTabCnt, er.rptTabCnt = t1, t2
	er.source, er.target = source, target

}

func (er *GGER) setDB2Alias(db2list []DB2List) {
	if er.program == "EXTRACT" {
		er.source = er.source + "(" + getDB2Info(er.source, db2list) + ")"
	}
	if er.program == "REPLICAT" {
		er.target = er.target + "(" + getDB2Info(er.target, db2list) + ")"
	}
}

func (er *GGER) CollectFileList(home string) ([]string, []string) {
	var res, resTemp []string
	if isExist(er.paramFile) {
		res = append(res, er.paramFile)
	}
	if len(er.obeyFiles) > 0 {
		res = append(res, er.obeyFiles...)
	}
	if er.propsFile != "" {
		res = append(res, er.propsFile)
	}
	if er.propertiesFile != "" {
		res = append(res, er.propertiesFile)
	}
	//report file
	filepath.Walk(filepath.Join(home, "dirrpt"), func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if strings.HasPrefix(f.Name(), er.group) && strings.HasSuffix(f.Name(), ".rpt") {
				res = append(res, filepath.Join(home, "dirrpt", f.Name()))
			}
		}
		return nil
	})

	//info detail, info showch , stats , send
	doWrite := func(wg *sync.WaitGroup, content, file string) {
		defer wg.Done()
		writeContent(content, file)
		if isExist(file) {
			res = append(res, file)
			resTemp = append(resTemp, file)
		}
	}
	var wg sync.WaitGroup
	wg.Add(5)
	go doWrite(&wg, er.getSendStatus(home), filepath.Join("/tmp", er.group+"_sendstatus.txt"))
	go doWrite(&wg, er.getInfo(home), filepath.Join("/tmp", er.group+"_info.txt"))
	go doWrite(&wg, er.getInfoShowch(home), filepath.Join("/tmp", er.group+"_showch.txt"))
	go doWrite(&wg, er.getInfoDetail(home), filepath.Join("/tmp", er.group+"_infodetail.txt"))
	go doWrite(&wg, er.getStats(home, "totalsonly *.*, reportrate min"), filepath.Join("/tmp", er.group+"_stats.txt"))
	wg.Wait()
	return res, resTemp
}
