package plugin

import (
	"database/sql"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/monitoring-system/dbtest/sqldiff"
)

var dataLoaderPluginMap map[string]func(*config.Config) interfaces.DataLoader = nil
var queryLoaderPluginMap map[string]func(*config.Config) interfaces.QueryLoader = nil
var comparePluginMap map[string]func(*config.Config) interfaces.SqlResultComparer = nil
var cellFilterMap map[string]func(conf *config.Config) interfaces.CellFilter = nil

func init() {
	dataLoaderPluginMap = make(map[string]func(*config.Config) interfaces.DataLoader)
	queryLoaderPluginMap = make(map[string]func(*config.Config) interfaces.QueryLoader)
	comparePluginMap = make(map[string]func(*config.Config) interfaces.SqlResultComparer)
	cellFilterMap = make(map[string]func(conf *config.Config) interfaces.CellFilter)

	dataLoaderPluginMap["dummy"] = newDummyDataLoader
	queryLoaderPluginMap["dummy"] = newDummyQueryLoader
	comparePluginMap["standard"] = newStandardCompare
	cellFilterMap["dummy"] = newCellFilter
}

func GetDataLoader(name string, config *config.Config) interfaces.DataLoader {
	f, ok := dataLoaderPluginMap[name]
	if ok {
		return f(config)
	}
	return nil
}

func GetQueryLoader(name string, config *config.Config) interfaces.QueryLoader {
	f, ok := queryLoaderPluginMap[name]
	if ok {
		return f(config)
	}
	return nil
}

func GetCompareLoader(name string, config *config.Config) interfaces.SqlResultComparer {
	f, ok := comparePluginMap[name]
	if ok {
		return f(config)
	}
	return nil
}

type dummyDataLoader struct {
}

func newDummyDataLoader(config *config.Config) interfaces.DataLoader {
	return &dummyDataLoader{}
}
func newDummyQueryLoader(config *config.Config) interfaces.QueryLoader {
	return &dummyDataLoader{}
}
func newStandardCompare(config *config.Config) interfaces.SqlResultComparer {
	return &sqldiff.StandardComparer{}
}
func newCellFilter(config *config.Config) interfaces.CellFilter {
	return &dummyDataLoader{}
}

func (data *dummyDataLoader) LoadData() []string {
	return []string{
		"explain select * from mysql.user",
	}
}

func (data *dummyDataLoader) LoadQuery() []string {
	return []string{
		"select * from mysql.user",
	}
}

func (data *dummyDataLoader) Filter([]byte, []byte, *sql.ColumnType) bool {
	return false
}
