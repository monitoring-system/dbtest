package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/api"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevel(),
		Encoding:         "json",
		ErrorOutputPaths: []string{"stdout"},
		OutputPaths:      []string{"stdout"},
	}

	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	_, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("init logger failed, err=%v", err))
	}

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
	log.Fatal("start server failed", zap.String("err", engine.Run("0.0.0.0:8080").Error()))

}
