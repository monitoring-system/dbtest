package impl

import (
	"strings"
)

const String = "string"

// file loader will loader data and queries from a file
type StringLoader struct {
	// sql split with ;
	SQLStr   string
}

func (loader *StringLoader) LoadSql(dbName string) []string {
	return loader.load(dbName)
}


func (loader *StringLoader) Name() string {
	return String
}

func (loader *StringLoader) load(dbName string) []string {
	origin := strings.Split(loader.SQLStr, ";")
	sql := make([]string, len(origin)+1)
	sql = append(sql, "use "+dbName+";")
	for _, st := range origin {
		if st != "" {
			sql = append(sql, st)
		}
	}
	return sql
}
