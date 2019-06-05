package main

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/plugin"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"

	"github.com/pingcap/log"
)

func main() {
	db1, _ := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)
	db2, _ := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)

	dataLoaders := strings.Split(config.GetConf().DataLoaders, ",")
	queryLoaders := strings.Split(config.GetConf().QueryLoaders, ",")
	compare := plugin.GetCompareLoader(config.GetConf().Comparor)

	round := 1
	for {
		log.Info("start to run test", zap.Int("round", round))

		for _, name := range dataLoaders {
			dataLoader := plugin.GetDataLoader(name)
			log.Info("get data loader from registry", zap.String("name", name))
			for _, sql := range dataLoader.LoadData() {
				log.Info("start execute sql", zap.String("sql", sql))
				r1, _ := sqldiff.GetQueryResult(db1, sql)
				r2, _ := sqldiff.GetQueryResult(db1, sql)
				fmt.Printf("%v\n%v\n", r1, r2)
			}
		}
		for _, name := range queryLoaders {
			queryLoader := plugin.GetQueryLoader(name)
			log.Info("get query loader from registry", zap.String("name", name))

			for _, query := range queryLoader.LoadQuery() {
				log.Info("execute query %s", zap.String("query", query))
				same, err1, err2 := compare.CompareQuery(db1, db2, query)
				fmt.Printf("%v\n%v\n%v\n", same, err1, err2)
			}
		}

		round++
		if round > config.GetConf().Loop {
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(config.GetConf().LoopInterval)))
		}
	}
}
