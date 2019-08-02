package executor

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/util"
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
		t.Fatalf("can not connect to db %v; %v", err1, err2)
	}

	initTestDatabase(config.GetConf().StandardDB, "test_dbtest")

	executor := &Executor{mysql: mysql, tidb: tidb, tests: make(map[string]*types.TestResult)}
	// dbtest_0_1
	tt := &types.Test{
		SqlLoader:"string",
		DataStr:"create table ff (f float, de decimal(40, 20) signed, key (f));insert into ff values(-0.2, 29.88);select ROUND(f) from ff",
		Loop:1,
	}

	executor.Submit(tt, true, true)
}
