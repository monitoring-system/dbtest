package impl

import (
	"encoding/json"
	"github.com/monitoring-system/dbtest/config"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type RandgenLoader struct {
	Yy      string
	Zz      string
	Queries int
}

func (this *RandgenLoader) LoadSql(dbname string) []string {
	return getLoadDataResponse(this, dbname).SQLs
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
		log.Warn("read randgen server resp error", zap.Error(err))
	}
	data := &LoadDataResponse{}
	json.Unmarshal(body, data)
	return data
}

type LoadDataResponse struct {
	SQLs    []string `json:"sql"`
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
