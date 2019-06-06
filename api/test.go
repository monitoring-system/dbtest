package api

import (
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/executor"
	"net/http"
)

type server struct {
	executor *executor.Executor
}

func NewServer(executor *executor.Executor) *server {
	return &server{executor: executor}
}

func (server *server) NewTest(c *gin.Context) {
	test := &executor.TestConfig{}
	if err := c.ShouldBind(test); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}
	server.executor.Submit(test)
	c.JSON(http.StatusOK, NewOKResponse())
}
