package impl

import (
	"strings"
)

// file loader will loader data and queries from a file
type StringLoader struct {
	SQLStr string
}

func (loader *StringLoader) LoadData(dbName string) []string {
	return loader.load(dbName)
}

func (loader *StringLoader) LoadQuery(dbName string) []string {
	return loader.load(dbName)
}

func (loader *StringLoader) Name() string {
	return "string"
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
