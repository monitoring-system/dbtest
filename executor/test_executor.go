package executor

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/filter"
	"github.com/monitoring-system/dbtest/parser"
	"github.com/monitoring-system/dbtest/sqldiff"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	golog "log"
	"os"
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
	executor := &Executor{mysql: mysql, tidb: tidb, tests: make(map[string]*types.TestResult)}
	executor.reScheduleUnfinishedTests()
	return executor
}

func (executor *Executor) reScheduleUnfinishedTests() {
	list, err := types.ListUnFinishedTestResult()
	if err != nil {
		log.Warn("reschedule failed", zap.Error(err))
	}
	if len(list) > 0 {
		for _, result := range list {
			test, err := types.GetTestById(result.TestID)
			if err != nil {
				log.Warn("unable to get test config, skip", zap.Int64("testId", result.TestID), zap.Int64("resusltId", result.ID), zap.Error(err))
				continue
			}
			go executor.run(test, result)
		}
	}
}

//submit a Test to the executor
func (executor *Executor) Submit(test *types.Test) (*types.Test, error) {
	executor.lock.Lock()
	executor.lock.Unlock()
	err := test.Persistent()
	if err != nil {
		return nil, err
	}
	result := &types.TestResult{TestID: test.ID, Name: test.GetName(), Status: types.TestStatusPending, Loop: 0, Start: time.Now().Unix()}
	err = result.Persistent()
	if err != nil {
		return nil, err
	}
	go executor.run(test, result)
	return test, nil
}

type testScope struct {
	test         *types.Test
	result       *types.TestResult
	loopResult   *types.LoopResult
	dbName       string
	db1          *sql.DB
	db2          *sql.DB
	ignoreTables *util.Set
	logger       *golog.Logger
}

func (executor *Executor) run(test *types.Test, result *types.TestResult) {
	round := 1
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Info("execute loop failed")
				}
			}()
			logger, file, err := getLogger(test, round, "log", golog.LstdFlags)
			if err != nil {
				loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Start: time.Now().Unix(), Status: types.TestStatusSkip}
				loopResult.Persistent()
				return
			}
			defer file.Close()

			loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Start: time.Now().Unix(), Status: types.TestStatusOK}
			result.Loop += 1
			result.Status = types.TestStatusRunning
			dbName := fmt.Sprintf("dbtest_%d_%d", test.ID, round)
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
			scope := &testScope{test: test, result: result, loopResult: loopResult, dbName: dbName, db1: db1, db2: db2, ignoreTables: util.NewSet(), logger: logger}

			log.Info("start to run test", zap.String("TestName", test.TestName), zap.Int64("TestId", test.ID), zap.Int("round", round), zap.String("dbName", dbName))
			logger.Println("dbName", dbName)
			data := executor.execDML(scope)
			query := executor.execQuery(scope)

			loopResult.End = time.Now().Unix()
			if err := loopResult.Persistent(); err != nil {
				log.Warn("insert loop result failed", zap.Error(err))
			}
			if loopResult.Status != types.TestStatusOK {
				result.FailedLoopCount++
				logger.Println("test case failed")
				persistentData(test, round, data.String())
				persistentQuery(test, round, query.String())
			} else {
				logger.Println("test case OK")
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

func (executor *Executor) execQuery(scope *testScope) *bytes.Buffer {
	compare := scope.test.GetComparor()
	var queryBuf = bytes.Buffer{}
	queryLoader := scope.test.GetQueryLoaders()
	scope.logger.Println("load queries", zap.Int64("testId", scope.test.ID), zap.String("name", queryLoader.Name()))
	for _, query := range queryLoader.LoadQuery(scope.dbName) {
		if query == "" || len(query) == 0 {
			continue
		}
		parsed, shouldIgnore := shouldSkipStatement(scope.logger, query, scope.ignoreTables)
		if shouldIgnore {
			continue
		}

		if parsed.Rewrite {
			query = parsed.NewSql
		}

		queryBuf.WriteString(query)
		queryBuf.WriteString(";\n")
		log.Info("execute query", zap.Int64("testId", scope.test.ID), zap.String("query", query))
		diff, err1, err2 := compare.CompareQuery(scope.db1, scope.db2, query)
		if putIgnoreTable(scope.logger, parsed, scope.ignoreTables, err1) {
			continue
		}
		if putIgnoreTable(scope.logger, parsed, scope.ignoreTables, err2) {
			continue
		}
		ignore := false
		if err2 != nil {
			ignore = filter.FilterError(err2.Error(), query)
		}
		log.Info("done", zap.String("diff", diff), zap.Bool("ignore", ignore), zap.Error(err1), zap.Error(err2))
		if diff != "" && !ignore {
			scope.loopResult.Status = types.TestStatusFail
			scope.logger.Println("compare sql result failed", query)
			scope.logger.Println(diff)
		}
	}
	return &queryBuf
}

func (executor *Executor) execDML(scope *testScope) *bytes.Buffer {
	var dataBUf = bytes.Buffer{}
	dataLoader := scope.test.GetDataLoaders()
	log.Info("using data loader to load data", zap.Int64("testId", scope.test.ID), zap.String("name", dataLoader.Name()))
	for _, statement := range dataLoader.LoadData(scope.dbName) {
		if statement == "" || len(statement) == 0 {
			continue
		}

		parsed, shouldIgnore := shouldSkipStatement(scope.logger, statement, scope.ignoreTables)
		if shouldIgnore {
			continue
		}

		if parsed.Rewrite {
			statement = parsed.NewSql
		}

		dataBUf.WriteString(statement)
		dataBUf.WriteString(";\n")
		log.Info("start execute statement", zap.Int64("testId", scope.test.ID), zap.String("statement", statement))
		r1, err1 := sqldiff.GetQueryResult(scope.db1, statement)
		if putIgnoreTable(scope.logger, parsed, scope.ignoreTables, err1) {
			continue
		}
		r2, err2 := sqldiff.GetQueryResult(scope.db2, statement)
		if putIgnoreTable(scope.logger, parsed, scope.ignoreTables, err2) {
			continue
		}
		if err2 != nil && !filter.FilterError(err2.Error(), statement) {
			scope.logger.Println(fmt.Sprintf("%v\n%v\n%v\n%v\n", r1, r2, err1, err2))
		}
	}
	return &dataBUf
}

func putIgnoreTable(logger *golog.Logger, parsed *parser.Result, ignored *util.Set, err error) bool {
	if config.GetConf().TraceAllErrors {
		return false
	}
	if err != nil && parsed.IsDDL {
		for _, parsedTableName := range parsed.TableName {
			_ = ignored.Put(parsedTableName)
		}
		return true
	}
	return false
}

func shouldSkipStatement(logger *golog.Logger, statement string, ignoreTables *util.Set) (*parser.Result, bool) {
	parsed, err := parser.Parse(statement)
	if config.GetConf().TraceAllErrors {
		return parsed, false
	}
	if err != nil {
		return parsed, true
	}

	if parsed.IgnoreSql {
		return parsed, true
	}

	shouldIgnore := false
	for _, parsedTableName := range parsed.TableName {
		if ignoreTables.Contains(parsedTableName) {
			log.Info("ignore failed table with failed ddl", zap.String("statement", statement), zap.String("table", parsedTableName))
			shouldIgnore = true
			break
		}
	}
	return parsed, shouldIgnore
}

func getLogger(test *types.Test, round int, suffix string, flag int) (*golog.Logger, *os.File, error) {
	logDir := fmt.Sprintf("results/logs/%d", test.ID)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%d.%s", logDir, round, suffix), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	logger := golog.New(f, "", flag)
	return logger, f, nil
}

func persistentData(test *types.Test, round int, data string) {
	dataLogger, dataFile, _ := getLogger(test, round, "sql", 0)
	defer dataFile.Close()
	dataLogger.Println(data)
}

func persistentQuery(test *types.Test, round int, query string) {
	dataLogger, dataFile, _ := getLogger(test, round, "query", 0)
	defer dataFile.Close()
	dataLogger.Println(query)
}
