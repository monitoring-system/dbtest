package cmd

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var StartCmd = &cobra.Command{
	Use:   "start [options ]",
	Short: "start the db test web server",
	Long:  "start the db test web server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func init() {
	config.Conf = &config.Config{}
	StartCmd.Flags().StringVar(&config.Conf.StandardDB, "standard-db", "root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	StartCmd.Flags().StringVar(&config.Conf.TestDB, "test-db", "root:@tcp(127.0.0.1:4000)/?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.Parse()
}

func StartServer() {
	initDatabase(config.GetConf().StandardDB)
	engine := gin.Default()

	MySQL, err1 := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)
	TiDB, err2 := util.OpenDBWithRetry("mysql", config.GetConf().TestDB)
	if err1 != nil || err2 != nil {
		log.Fatal("can not connect to db", zap.Error(err1), zap.Error(err2))
	}

	server := api.NewServer(executor.New(MySQL, TiDB))

	engine.POST("/tests", server.NewTest)
	engine.GET("/tests", server.ListTest)
	engine.GET("/tests/:id", server.GetTest)
	engine.GET("/results", server.ListTestResult)
	engine.GET("/results/:id/detail", server.ListLoopResult)

	engine.POST("/addfilter", server.AddFilter)

	log.Fatal("StartServer server failed", zap.String("err", engine.Run("0.0.0.0:8080").Error()))

}

var models = []interface{}{&types.Test{}, &types.TestResult{}, &types.LoopResult{}}

func initDatabase(dsn string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("can not parse the mysql configuration")
	}
	if cfg.DBName == "" {
		db, err := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
		if err != nil {
			panic("can not connect to database")
		}
		db.Exec("create database if not exists dbtest")
		db.Close()
		cfg.DBName = "dbtest"
	}
	db.InitDatabase(cfg.FormatDSN(), models)
}
