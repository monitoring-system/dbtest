package executor

import (
	"bytes"
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

//submit a Test to the executor
func (executor *Executor) Submit(test *types.Test) (*types.Test, error) {
	executor.lock.Lock()
	executor.lock.Unlock()
	err := test.Persistent()
	if err != nil {
		return nil, err
	}
	result := &types.TestResult{TestID: test.ID, Name: test.Name(), Status: types.TestStatusPending, Loop: 0, Start: time.Now().Unix()}
	err = result.Persistent()
	if err != nil {
		return nil, err
	}
	go executor.run(test, result)
	return test, nil
}

type testScope struct {
	test       *types.Test
	result     *types.TestResult
	loopResult *types.LoopResult
	dbName     string
	db1        *sql.DB
	db2        *sql.DB
}

func (executor *Executor) run(test *types.Test, result *types.TestResult) {
	round := 1
	for {
		func() {
			loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Start: time.Now().Unix(), Status: types.TestStatusOK}
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

			scope := &testScope{test: test, result: result, loopResult: loopResult, dbName: dbName, db1: db1, db2: db2}
			executor.execDML(scope)
			executor.execQuery(scope)

			loopResult.End = time.Now().Unix()
			if err := loopResult.Persistent(); err != nil {
				log.Warn("insert loop result failed", zap.Error(err))
			}
			if loopResult.FailedDML != "" || loopResult.FailedQuery != "" {
				result.FailedLoopCount++
			}
			if err := result.Update(); err != nil {
				log.Warn("update result failed", zap.Error(err))
			}
		}()

		round++
		if round > test.GetLoop() {
			result.Status = types.TestStatusDone
			result.End = time.Now().Unix()
			if err := result.Update(); err != nil {
				log.Warn("update result failed", zap.Error(err))
			}
			break
		} else {
			time.Sleep(time.Duration(time.Second * time.Duration(test.GetLoopInterval())))
		}
	}
}

func (executor *Executor) execQuery(scope *testScope) {
	compare := scope.test.GetComparor()
	var queryBuf = bytes.Buffer{}
	var failedQueryBuf = bytes.Buffer{}
	for _, queryLoader := range scope.test.GetQueryLoaders() {
		log.Info("load queries", zap.Int64("testId", scope.test.ID), zap.String("name", queryLoader.Name()))
		for _, query := range queryLoader.LoadQuery(scope.dbName) {
			if query == "" || len(query) == 0 {
				continue
			}
			queryBuf.WriteString(query)
			queryBuf.WriteString(";")
			log.Info("execute query", zap.Int64("testId", scope.test.ID), zap.String("query", query))
			same, err1, err2 := compare.CompareQuery(scope.db1, scope.db2, query)
			ignore := false
			if err2 != nil {
				ignore = filter.FilterError(err2.Error(), query)
			}
			fmt.Printf("%v\n%v\n%v\n%v\n", same, err1, err2, ignore)
			if !same && !ignore {
				failedQueryBuf.WriteString(query)
				failedQueryBuf.WriteString(";")
				scope.loopResult.Status = types.TestStatusFail
			}
		}
	}
	scope.loopResult.Query = queryBuf.String()
	scope.loopResult.FailedQuery = failedQueryBuf.String()
}

func (executor *Executor) execDML(scope *testScope) {
	var dataBUf = bytes.Buffer{}
	var failedDML = bytes.Buffer{}
	for _, dataLoader := range scope.test.GetDataLoaders() {
		log.Info("using data loader to load data", zap.Int64("testId", scope.test.ID), zap.String("name", dataLoader.Name()))
		for _, statement := range dataLoader.LoadData(scope.dbName) {
			if statement == "" || len(statement) == 0 {
				continue
			}
			dataBUf.WriteString(statement)
			dataBUf.WriteString(";")
			log.Info("start execute statement", zap.Int64("testId", scope.test.ID), zap.String("statement", statement))
			r1, err1 := sqldiff.GetQueryResult(scope.db1, statement)
			r2, err2 := sqldiff.GetQueryResult(scope.db2, statement)
			fmt.Printf("%v\n%v\n%v\n%v\n", r1, r2, err1, err2)
			if err2 != nil && !filter.FilterError(err2.Error(), statement) {
				failedDML.WriteString(statement)
			}
		}
	}
	scope.loopResult.DML = dataBUf.String()
	scope.loopResult.FailedDML = failedDML.String()
}
