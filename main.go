package main

import (
	"flag"
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

var conf *config.Config

func init() {
	conf = &config.Config{}
	flag.IntVar(&conf.Loop, "loop", 1, "the loop count")
	flag.IntVar(&conf.LoopInterval, "loop-interval", 10, "the second to sleep after a loop is finished")
	flag.StringVar(&conf.DataLoaders, "data-loaders", "dummy", "a list of data loader names split by comma")
	flag.StringVar(&conf.QueryLoaders, "query-loaders", "dummy", "a list of query loader names split by comma")
	flag.StringVar(&conf.Comparor, "comparor", "standard", "the compare plugin")
	flag.StringVar(&conf.Comparor, "cell-filter", "standard", "the cell filter plugin")
	flag.StringVar(&conf.StandardDB, "standard-db", "root:@tcp(127.0.0.1:3306)/tep?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.StringVar(&conf.TestDB, "test-db", "root:@tcp(127.0.0.1:4000)/tep?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.Parse()
}

func main() {
	db1, _ := util.OpenDBWithRetry("mysql", conf.StandardDB)
	db2, _ := util.OpenDBWithRetry("mysql", conf.TestDB)

	dataLoaders := strings.Split(conf.DataLoaders, ",")
	queryLoaders := strings.Split(conf.QueryLoaders, ",")
	compare := plugin.GetCompareLoader(conf.Comparor, conf)

	round := 1
	for {
		log.Info("start to run test", zap.Int("round", round))

		for _, name := range dataLoaders {
			dataLoader := plugin.GetDataLoader(name, conf)
			log.Info("get data loader from registry", zap.String("name", name))
			for _, sql := range dataLoader.LoadData() {
				log.Info("start execute sql", zap.String("sql", sql))
				r1, _ := sqldiff.GetQueryResult(db1, sql)
				r2, _ := sqldiff.GetQueryResult(db1, sql)
				fmt.Printf("%v\n%v\n", r1, r2)
			}
		}
		for _, name := range queryLoaders {
			queryLoader := plugin.GetQueryLoader(name, conf)
			log.Info("get query loader from registry", zap.String("name", name))

			for _, query := range queryLoader.LoadQuery() {
				log.Info("execute query %s", zap.String("query", query))
				same, err1, err2 := compare.CompareQuery(db1, db2, query)
				fmt.Printf("%v\n%v\n%v\n", same, err1, err2)
			}
		}

		round++
		if round > conf.Loop {
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(conf.LoopInterval)))
		}
	}
}
