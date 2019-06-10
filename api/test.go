package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/monitoring-system/dbtest/filter"
	"net/http"
	"os"
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

func (server *server) AddFilter(c *gin.Context) {
	h, err := c.FormFile("key")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}

	os.Mkdir(filter.FilterPATH, os.ModePerm)
	filename := fmt.Sprintf("%s/%s", filter.FilterPATH, h.Filename)

	if err := c.SaveUploadedFile(h, filename); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	if err := filter.AddFilter(filename); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"filepath": "http://127.0.0.1:8000/" + filename})
}
