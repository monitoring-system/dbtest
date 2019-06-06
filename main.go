package main

import (
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/api"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func main() {
	engine := gin.Default()
	exec := &executor.Executor{}
	server := api.NewServer(exec)

	engine.POST("/test", server.NewTest)
	log.Fatal("start server failed", zap.String("err", engine.Run("0.0.0.0:8080").Error()))
}
