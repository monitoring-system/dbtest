package impl

import (
	"io/ioutil"
	"strings"
)

// file loader will loader data and queries from a file
type FileDataLoader struct {
	FileName string
}

func (loader *FileDataLoader) LoadData(dbName string) []string {
	return loader.load(dbName)
}

func (loader *FileDataLoader) LoadQuery(dbName string) []string {
	return loader.load(dbName)
}

func (loader *FileDataLoader) Name() string {
	return "file"
}

func (loader *FileDataLoader) load(dbName string) []string {
	sqlBytes, err := ioutil.ReadFile(loader.FileName)
	if err != nil {
		return nil
	}

	origin := strings.Split(string(sqlBytes), ";\n")
	sql := make([]string, len(origin)+1)
	sql = append(sql, "use "+dbName+";")
	for _, st := range origin {
		sql = append(sql, st)
	}
	return sql
}
