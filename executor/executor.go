package executor

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/filter"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"github.com/dqinyuan/sqlparser"
	"go.uber.org/zap"
	golog "log"
	"os"
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
			go executor.run(test, result, false)
		}
	}
}

//submit a Test to the executor
// db will not be dropped when test over, if dbRetain set True
func (executor *Executor) Submit(test *types.Test, syn bool, dbRetain bool) (*types.Test, error) {
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
	if syn {
		executor.run(test, result, dbRetain)
	} else {
		go executor.run(test, result, dbRetain)
	}
	return test, nil
}

type testScope struct {
	test         *types.Test
	result       *types.TestResult
	loopResult   *types.LoopResult
	dbName       string
	// mysql
	db1          *sql.DB
	// tidb
	db2          *sql.DB
	logger       *golog.Logger
}

func (executor *Executor) run(test *types.Test, result *types.TestResult, dbRetain bool) {
	round := 1
	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Info("execute loop failed", zap.String("err", fmt.Sprintf("%v", err)))
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
			log.Info("db info", zap.String("dbname", dbName))
			_, err = executor.mysql.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			if err != nil {
				log.Info("fail to create database in mysql", zap.Error(err))
				loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Start: time.Now().Unix(), Status: types.TestStatusSkip}
				loopResult.Persistent()
				return
			}
			_, err = executor.tidb.Exec("CREATE DATABASE IF NOT EXISTS  " + dbName)
			if err != nil {
				log.Info("fail to create database in tidb", zap.Error(err))
				loopResult := &types.LoopResult{TestID: test.ID, Loop: round, Start: time.Now().Unix(), Status: types.TestStatusSkip}
				loopResult.Persistent()
				return
			}
			if !dbRetain {
				defer executor.mysql.Exec("DROP DATABASE IF EXISTS  " + dbName)
				defer executor.tidb.Exec("DROP DATABASE IF EXISTS  " + dbName)
			}

			cfg, _ := mysql.ParseDSN(config.GetConf().StandardDB)
			cfg.DBName = dbName
			db1, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			cfg, _ = mysql.ParseDSN(config.GetConf().TestDB)
			cfg.DBName = dbName
			db2, _ := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
			defer db1.Close()
			defer db2.Close()
			scope := &testScope{test: test, result: result, loopResult: loopResult, dbName: dbName, db1: db1, db2: db2, logger: logger}

			log.Info("start to run test",
				zap.String("TestName", test.TestName),
				zap.Int64("TestId", test.ID),
				zap.Int("round", round),
				zap.String("dbName", dbName))
			logger.Println("dbName", dbName)
			// exec
			data := executor.execSql(scope, scope.test.GetComparor())

			loopResult.End = time.Now().Unix()
			if err := loopResult.Persistent(); err != nil {
				log.Warn("insert loop result failed", zap.Error(err))
			}
			if loopResult.Status != types.TestStatusOK {
				result.FailedLoopCount++
				logger.Println("test case failed")
				persistentData(test, round, data.String())
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

func (executor *Executor) execSql(scope *testScope, compare interfaces.SqlResultComparer) *bytes.Buffer  {

	// store all sql in Buf
	sqlBuf := &bytes.Buffer{}
	loader := scope.test.GetSqlLoader()

	scope.logger.Println("load sqls ", fmt.Sprintf("testId=%d", scope.test.ID))

	sqls := loader.LoadSql(scope.dbName)
	if len(sqls) == 0 {
		scope.logger.Println("no sql is generated")
		log.Warn("no sql is found")
		return sqlBuf
	}
	execCount := 0

	replacer := strings.NewReplacer("\n", "", "\t", "")

	for _, sql := range sqls {
		if sql == "" {
			continue
		}

		ast, err := sqlparser.ParseStrictDDL(sql)
		if err != nil {
			// if run in there, perhaps we need to replace current sql parser
			scope.logger.Printf("sql parse error %v, sql is %s\n", err, sql)
			log.Warn("sql parse error", zap.String("sql", sql), zap.Error(err))
		} else {
			// filter
			newAst, comment := filter.FiltByCtx(ast)
			if newAst == nil {
				scope.logger.Printf("sql filtered: \n    %s,\n    err Type:%s\n", sql, comment)
				continue
			}
			// rewrite sql
			sql = replacer.Replace(sqlparser.String(newAst))
		}

		// log sql after rewrite
		sqlBuf.WriteString(sql)
		sqlBuf.WriteString(";\n")
		execCount++

		// diff
		diff, err1, err2 := compare.CompareQuery(scope.db1, scope.db2, sql)
		log.Info("done", zap.String("diff", diff), zap.Error(err1), zap.Error(err2))

		if diff != "" {
			scope.loopResult.Status = types.TestStatusFail
			scope.logger.Println("compare sql result failed", sql)
			scope.logger.Println(diff)
		}

		// update filter context by result
		filter.UpdateCtxByExecResult(ast, err1, err2)
	}

	filter.ClearFilterCtx()
	scope.logger.Println("executed query count", execCount)
	return sqlBuf
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
