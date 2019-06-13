package executor

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/filter"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type Executor struct {
	mysql *sql.DB
	tidb  *sql.DB
	tests map[string]*types.TestResult
	lock  sync.Mutex
}

func New(mysql, tidb *sql.DB) *Executor {
	return &Executor{mysql: mysql, tidb: tidb, tests: make(map[string]*types.TestResult)}
}

func (executor *Executor) Submit(test *types.Test) (*types.Test, error) {
	executor.lock.Lock()
	executor.lock.Unlock()
	err := test.Persistent()
	if err != nil {
		return nil, err
	}
	result := &types.TestResult{TestID: test.ID, Name: test.Name(), Status: types.TestStatusPending, Loop: 0}
	err = result.Persistent()
	if err != nil {
		return nil, err
	}
	go executor.run(test, result)
	return test, nil
}

func (executor *Executor) run(test *types.Test, result *types.TestResult) {
	compare := test.GetComparor()
	round := 1
	for {
		func() {
			result.Loop += 1
			result.Status = types.TestStatusRunning
			dbName := "tbl" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")
			executor.mysql.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			executor.tidb.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			defer executor.mysql.Exec("DROP DATABASE IF EXISTS  " + dbName)
			defer executor.tidb.Exec("DROP DATABASE IF EXISTS  " + dbName)

			cfg, _ := mysql.ParseDSN(config.GetConf().StandardDB)
			cfg.DBName = dbName
			db1, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			cfg, _ = mysql.ParseDSN(config.GetConf().TestDB)
			cfg.DBName = dbName
			db2, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			defer db1.Close()
			defer db2.Close()

			log.Info("start to run test", zap.Int("round", round))

			//set := util.Set{}
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
					ignore := false
					if err2 != nil {
						ignore = filter.FilterError(err2.Error(), query)
					}
					fmt.Printf("%v\n%v\n%v\n%v\n", same, err1, err2, ignore)
				}
			}
			loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Detail: ""}
			if err := loopResult.Persistent(); err != nil {
				log.Warn("insert loop result failed", zap.Error(err))
			}
			if err := result.Update(); err != nil {
				log.Warn("update result failed", zap.Error(err))
			}
		}()

		round++
		if round > test.GetLoop() {
			result.Status = types.TestStatusDone
			if err := result.Update(); err != nil {
				log.Warn("update result failed", zap.Error(err))
			}
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(test.GetLoopInterval())))
		}
	}
}
