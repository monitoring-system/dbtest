package executor

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/util"
	"go.uber.org/zap"
	"log"
	"testing"
)

var models = []interface{}{&types.Test{}, &types.TestResult{}, &types.LoopResult{}}

func initTestDatabase(dsn string, dbname string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("can not parse the mysql configuration")
	}
	if cfg.DBName == "" {
		db, err := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
		if err != nil {
			panic("can not connect to database")
		}
		db.Exec(fmt.Sprintf("create database if not exists %s", dbname))
		db.Close()
		cfg.DBName = dbname
	}
	db.InitDatabase(cfg.FormatDSN(), models)
}


func TestExecutor_Submit(t *testing.T) {
	t.SkipNow()
	config.Conf = &config.Config{
		StandardDB:"root:123456@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local",
		TestDB:"root:@tcp(127.0.0.1:4000)/?charset=utf8&parseTime=True&loc=Local",
	}

	mysql, err1 := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)
	tidb, err2 := util.OpenDBWithRetry("mysql", config.GetConf().TestDB)
	if err1 != nil || err2 != nil {
		log.Fatal("can not connect to db", zap.Error(err1), zap.Error(err2))
	}

	initTestDatabase(config.GetConf().StandardDB, "test_dbtest")

	executor := &Executor{mysql: mysql, tidb: tidb, tests: make(map[string]*types.TestResult)}
	// dbtest_0_1
	tt := &types.Test{
		SqlLoader:"string",
		DataStr:"create table t (a float);insert into t values(0.12345678910111213);select * from t",
		Loop:1,
	}

	executor.Submit(tt, true, true)
}


/*func TestExecSql(t *testing.T) {
	sqls := []string{
		"drop table if exists nums",
		"create table nums (num int)",
		"insert into nums (num) values(10)",
		"select * from nums",
	}

	test := &types.Test{
		SqlLoader:"string",
		DataStr:strings.Join(sqls, ";"),
	}

	db1, _ := sql.Open("ramsql", "mysql")
	db2, _ := sql.Open("ramsql", "tidb")

	scope := &testScope{test:test, loopResult:&types.LoopResult{}, logger: golog.New(os.Stdout,
		"", golog.LstdFlags), dbName:"test", db1: db1, db2: db2}

	exec := &Executor{}
	eSqls := exec.execSql(scope)
	fmt.Println(eSqls.String())
}
*/