package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/api"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zap.Config{
		Level: zap.NewAtomicLevel(),
		Encoding: "json",
		ErrorOutputPaths: []string{"stdout"},
		OutputPaths: []string{"stdout"},
	}

	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	_, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("init logger failed, err=%v", err))
	}


	engine := gin.Default()
	exec := &executor.Executor{}
	server := api.NewServer(exec)

	engine.POST("/test", server.NewTest)
	engine.POST("/addfilter", server.AddFilter)
	log.Fatal("start server failed", zap.String("err", engine.Run("0.0.0.0:8080").Error()))

}


