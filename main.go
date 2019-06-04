package main

import (
	"fmt"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/plugin"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"
)

func main() {
	db1, _ := util.OpenDBWithRetry("mysql", "root:@tcp(127.0.0.1:3306)/tep?charset=utf8&parseTime=True&loc=Local")
	db2, _ := util.OpenDBWithRetry("mysql", "root:@tcp(127.0.0.1:3306)/tep?charset=utf8&parseTime=True&loc=Local")

	//
	conf := &config.Config{}

	dataLoader := plugin.GetDataLoader("dummy", conf)
	queryLoad := plugin.GetQueryLoader("dummy", conf)

	compare := plugin.GetCompareLoader("standard", conf)

	for _, sql := range dataLoader.LoadData() {
		fmt.Println("execute %s\n", sql)
		r1, _ := sqldiff.GetQueryResult(db1, sql)
		r2, _ := sqldiff.GetQueryResult(db1, sql)
		fmt.Printf("%v\n%v\n", r1, r2)
	}

	for _, sql := range queryLoad.LoadQuery() {
		fmt.Println("execute query %s\n", sql)
		same, err1, err2 := compare.CompareQuery(db1, db2, sql)
		fmt.Printf("%v\n%v\n%v\n", same, err1, err2)
	}
}
