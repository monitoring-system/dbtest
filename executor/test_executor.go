package executor

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/interfaces"
	"strings"
	"time"

	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"

	"github.com/pingcap/log"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type Executor struct {
	MySQL *sql.DB
	TiDB  *sql.DB
}

type TestConfig interface {
	GetDataLoaders() []interfaces.DataLoader
	GetQueryLoaders() []interfaces.QueryLoader
	GetComparor() interfaces.SqlResultComparer
	GetLoop() int
	GetLoopInterval() int
}

func (executor *Executor) Submit(test TestConfig) {
	go executor.run(test)
}

func (executor *Executor) run(test TestConfig) {

	compare := test.GetComparor()
	round := 1
	for {
		func() {
			dbName := "tbl" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")
			executor.MySQL.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			executor.TiDB.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			defer executor.MySQL.Exec("DROP DATABASE IF EXISTS  " + dbName)
			defer executor.TiDB.Exec("DROP DATABASE IF EXISTS  " + dbName)

			cfg, _ := mysql.ParseDSN(config.GetConf().StandardDB)
			cfg.DBName = dbName
			db1, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			cfg, _ = mysql.ParseDSN(config.GetConf().TestDB)
			cfg.DBName = dbName
			db2, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			defer db1.Close()
			defer db2.Close()

			log.Info("start to run test", zap.Int("round", round))

			for _, dataLoader := range test.GetDataLoaders() {
				log.Info("using data loader to load data", zap.String("name", dataLoader.Name()))
				for _, statement := range dataLoader.LoadData(dbName) {
					if statement == "" || len(statement) == 0 {
						continue
					}
					log.Info("start execute statement", zap.String("statement", statement))
					r1, err1 := sqldiff.GetQueryResult(db1, statement)
					r2, err2 := sqldiff.GetQueryResult(db2, statement)
					fmt.Printf("%v\n%v\n%v\n%v\n", r1, r2, err1, err2)
				}
			}
			for _, queryLoader := range test.GetQueryLoaders() {
				log.Info("load queries", zap.String("name", queryLoader.Name()))
				for _, query := range queryLoader.LoadQuery(dbName) {
					if query == "" || len(query) == 0 {
						continue
					}
					log.Info("execute query", zap.String("query", query))
					same, err1, err2 := compare.CompareQuery(db1, db2, query)
					fmt.Printf("%v\n%v\n%v\n", same, err1, err2)
				}
			}
		}()

		round++
		if round > test.GetLoop() {
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(test.GetLoopInterval())))
		}
	}
}
