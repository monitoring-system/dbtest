package impl

import (
	"encoding/json"
	"github.com/monitoring-system/dbtest/config"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type RandgenLoader struct {
	Yy      string
	Zz      string
	Queries int

	locker sync.RWMutex
	data   []string
	query  []string
}

func (test *RandgenLoader) LoadData(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	response := getLoadDataResponse(test, db)
	if response == nil {
		test.data = nil
		test.query = nil
		return nil
	}
	test.data = response.SQLs
	test.query = response.Queries
	return test.data
}

func (test *RandgenLoader) LoadQuery(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	return test.query
}

func (test *RandgenLoader) Name() string {
	return "randgen"
}

func getLoadDataResponse(randgen *RandgenLoader, db string) *LoadDataResponse {
	payload := &LoadDataRequest{Yy: randgen.Yy, ZZ: randgen.Zz, DB: db, Queries: randgen.Queries}
	resp, err := http.Post(config.Conf.RandGenServer, "application/json",
		strings.NewReader(getLoadDataRequestString(payload)))
	if err != nil || resp == nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	data := &LoadDataResponse{}
	json.Unmarshal([]byte(body), data)
	return data
}

type LoadDataResponse struct {
	SQLs    []string `json:"sql"`
	Queries []string `json:"queries"`
}

func getLoadDataRequestString(payload *LoadDataRequest) string {
	bytes, _ := json.Marshal(payload)
	return string(bytes)
}

type LoadDataRequest struct {
	Yy      string `json:"yy"`
	ZZ      string `json:"zz"`
	DB      string `json:"db"`
	Queries int    `json:"queries"`
}
