package executor

import (
	"fmt"
	"strings"
	"time"

	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/plugin"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"

	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type Executor struct {
}

type TestConfig struct {
	DataLoaders  string `json:"dataLoaders,omitempty"`
	QueryLoaders string `json:"queryLoaders,omitempty"`
	Comparor     string `json:"comparor,omitempty"`

	Loop         int `json:"loop,omitempty"`
	LoopInterval int `json:"loopInterval,omitempty"`
}

func (executor *Executor) Submit(test *TestConfig) {
	go executor.run(test)
}

func (executor *Executor) run(test *TestConfig) {
	db1, _ := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)
	db2, _ := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)

	dataLoaders := strings.Split(test.DataLoaders, ",")
	queryLoaders := strings.Split(test.QueryLoaders, ",")
	compare := plugin.GetCompareLoader(test.Comparor)

	round := 1
	for {
		log.Info("start to run test", zap.Int("round", round))

		for _, name := range dataLoaders {
			dataLoader := plugin.GetDataLoader(name)
			if dataLoader == nil {
				log.Warn("can not found the data loader from registry", zap.String("name", name))
				continue
			}
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
			if queryLoader == nil {
				log.Warn("can not found the query loader from registry", zap.String("name", name))
				continue
			}
			log.Info("get query loader from registry", zap.String("name", name))

			for _, query := range queryLoader.LoadQuery() {
				log.Info("execute query %s", zap.String("query", query))
				same, err1, err2 := compare.CompareQuery(db1, db2, query)
				fmt.Printf("%v\n%v\n%v\n", same, err1, err2)
			}
		}

		round++
		if round > test.Loop {
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(test.LoopInterval)))
		}
	}
}
