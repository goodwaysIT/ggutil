package ogg

import (
	"strings"
)

// GGSoft 定义了 OGG 软件安装的相关信息
type GGSoft struct {
	home    string // OGG 安装主目录路径
	v       string // OGG 版本简要字符串 (例如 output of 'dblogin useridalias <alias> version')
	ggType  string // OGG 产品类型 (例如: Oracle, BigData, MySQL - typically from 'info all' or version string)
	oggFor  string // OGG 针对的具体数据库或应用 (例如: DB2, PostgreSQL - often part of version string)
	version string // OGG 详细版本号 (e.g., Version 19.1.0.0.230517 OGGCORE_19.1.0.0.0_PLATFORMS_230415.0100_FBO)
}

// GGInst 定义了一个 OGG 实例的结构
type GGInst struct {
	soft       *GGSoft  // 指向该实例所属的 GGSoft
	Er         []GGER   // 该实例下的 Extract 和 Replicat 进程列表
	Mgr        GGMgr    // 该实例下的 Manager 进程
	infoall    string   // 'info all' 命令的原始输出, 用于后续解析
	globalFile string   // GLOBALS 文件路径 (e.g., <OGG_HOME>/GLOBALS)
	errLog     string   // ggserr.log 文件路径 (e.g., <OGG_HOME>/ggserr.log)
	jarFiles   []string // JAR 文件列表 (例如 OGG for BigData, from dir <OGG_HOME>/lib/java/*)
}

// GGPgm 定义了 OGG 进程的通用属性
type GGPgm struct {
	program   string   // 进程名 (例如 EORA, RORA, MGR)
	status    string   // 进程状态 (例如 RUNNING, STOPPED, ABENDED, STARTING)
	rptFile   string   // 报告文件路径 (e.g., <OGG_HOME>/dirrpt/<PROGRAM>.rpt)
	paramFile string   // 参数文件路径 (e.g., <OGG_HOME>/dirprm/<PROGRAM>.prm)
	obeyFiles []string // Obey 文件列表 (e.g., <OGG_HOME>/diroby/<PROGRAM>*.oby)
}

// GGER 定义了 Extract 或 Replicat 进程特有的属性
type GGER struct {
	GGPgm                            // 嵌入 GGPgm 通用属性
	group                     string // 进程组名 (通常与 program 相同)
	lag                       string // 延迟信息 (e.g., "00:00:02", "At EOF")
	chkptTime                 string // 检查点时间 (e.g., "2023-10-26 14:30:05")
	propsFile                 string // Java属性文件路径 (例如 OGG for BigData, e.g., <OGG_HOME>/dirprm/<PROGRAM>.properties)
	propertiesFile            string // (备用) 属性文件路径
	// config - 从参数文件或报告文件解析的配置信息
	prmTabCnt                 int    // 参数文件中定义的表数量 (解析 TABLE/MAP 语句)
	rptTabCnt                 int    // 报告文件中记录的已处理表数量
	source                    string // 源端描述 (例如 schema.table)
	target                    string // 目标端描述 (例如 schema.table)
}

// GGMgr 定义了 Manager 进程特有的属性
type GGMgr struct {
	GGPgm // 嵌入 GGPgm 通用属性
	// Manager specific fields can be added here, e.g., port number, auto-restart settings
}

// ParseProcessNamesFromInfoAll extracts process names (MGR, Extract groups, Replicat groups)
// from the output of an 'info all' ggsci command.
func ParseProcessNamesFromInfoAll(infoAllOutput string) []string {
	var processNames []string
	lines := strings.Split(infoAllOutput, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		parts := strings.Fields(trimmedLine) // Splits by one or more spaces
		if len(parts) == 0 {
			continue
		}

		programType := strings.ToUpper(parts[0])

		switch programType {
		case "MANAGER":
			if len(parts) >= 1 { // e.g., "MANAGER   RUNNING"
				processNames = append(processNames, "MGR") // 'view param MGR'
			}
		case "EXTRACT", "REPLICAT":
			if len(parts) >= 2 { // e.g., "EXTRACT   EXT1      RUNNING"
				processNames = append(processNames, parts[1]) // 'view param EXT1'
			}
		}
	}
	return UniqueStrings(processNames)
}

// UniqueStrings helper function to remove duplicates from a slice of strings.
func UniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
