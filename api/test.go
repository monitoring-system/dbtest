package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/monitoring-system/dbtest/filter"
	"net/http"
	"strconv"
)

type server struct {
	executor *executor.Executor
}

func NewServer(executor *executor.Executor) *server {
	return &server{executor: executor}
}

func (server *server) NewTest(c *gin.Context) {
	test := &types.Test{}
	if err := c.ShouldBind(test); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	}
	result, err := server.executor.Submit(test)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (server *server) ListTest(c *gin.Context) {
	list, err := types.ListTest()
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func (server *server) GetTest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	} else {
		result, err := types.GetTestById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, NewErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (server *server) ListTestResult(c *gin.Context) {
	list, err := types.ListResult()
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func (server *server) ListLoopResult(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
	} else {
		result, err := types.ListLoopResult(id)
		if err != nil {
			c.JSON(http.StatusNotFound, NewErrorResponse(err.Error()))
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (server *server) AddFilter(c *gin.Context) {
	h, err := c.FormFile("key")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}

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
