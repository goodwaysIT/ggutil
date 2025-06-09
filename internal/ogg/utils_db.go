package ogg

import "strings"

// DB2List holds information about a DB2 database and its network node.
type DB2List struct {
	alias       string // Database alias
	name        string // Database name
	node        string // Node name
	hostname    string // Hostname for the node
	serviceName string // Service name for the node
}

func getDB2List() []DB2List {
	listNodeCmd := "db2 list node directory"
	listDbcmd := "db2 list database directory"
	nodeRes, _ := execShell(listNodeCmd)
	dbRes, _ := execShell(listDbcmd)

	var list1, list2, dbList []DB2List
	var db2db, db2node, db2var DB2List
	var node, hostname, serviceName, alias, name string
	// db list node
	nodeLines := strings.Split(nodeRes, "\n")
	for _, v := range nodeLines {
		if strings.Contains(v, "Node name") {
			node = strings.Split(v, "=")[1]
			continue
		}
		if strings.Contains(v, "Hostname") {
			hostname = strings.Split(v, "=")[1]
			continue
		}
		if strings.Contains(v, "Service name") {
			serviceName = strings.Split(v, "=")[1]
			db2node.node, db2node.hostname, db2node.serviceName =
				strings.TrimSpace(node), strings.TrimSpace(hostname), strings.TrimSpace(serviceName)
			list1 = append(list1, db2node)
			node, hostname, serviceName = "", "", ""
		}
	}
	//db list db
	dbLines := strings.Split(dbRes, "\n")
	for _, v := range dbLines {
		if strings.Contains(v, "Database alias") {
			alias = strings.Split(v, "=")[1]
			continue
		}
		if strings.Contains(v, "Database name") {
			name = strings.Split(v, "=")[1]
			continue
		}
		if strings.Contains(v, "Node name") {
			node = strings.Split(v, "=")[1]
			db2db.node, db2db.name, db2db.alias =
				strings.TrimSpace(node), strings.TrimSpace(name), strings.TrimSpace(alias)
			list2 = append(list2, db2db)
			node, name, alias = "", "", ""
		}
	}
	for _, v2 := range list2 {
		db2var.name, db2var.alias, db2var.node = v2.name, v2.alias, v2.node
		for _, v1 := range list1 {
			if v1.node == v2.node {
				db2var.hostname, db2var.serviceName = v1.hostname, v1.serviceName
				break
			}
		}
		dbList = append(dbList, db2var)
	}
	return dbList
}

func getDB2Info(alias string, dbList []DB2List) string {
	var res string
	for _, v := range dbList {
		if v.alias == strings.ToUpper(alias) {
			res = v.hostname + ":" + v.serviceName + "/" + v.name
			break
		}
	}
	return res
}

func getOracleInfo(alias string) string {
	var res string
	dbRes, _ := execShell("tnsping " + alias)
	lines := strings.Split(dbRes, "\n")
	for _, v := range lines {
		if strings.Contains(v, "DESCRIPTION") {
			lines1 := strings.Split(v, ")")
			for _, v1 := range lines1 {
				if strings.Contains(v1, "HOST") {
					res = strings.TrimSpace(strings.Split(v1, "=")[1])
				}
				if strings.Contains(v1, "PORT") {
					res = res + ":" + strings.TrimSpace(strings.Split(v1, "=")[1])
				}
				if strings.Contains(v1, "SERVICE_NAME") {
					res = res + "/" + strings.TrimSpace(strings.Split(v1, "=")[1])
				}
			}
		}
	}
	return res
}
