package randgen

import (
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/monitoring-system/dbtest/sqldiff"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type RandGen struct {
	Yy           string `json:"yy"`
	ZZ           string `json:"zz"`
	Queries      int    `json:"queries"`
	Loop         int    `json:"loop"`
	LoopInterval int    `json:"loopInterval"`

	locker sync.Mutex
	data   []string
	query  []string
	first  bool
}

func (rg *RandGen) LoadData(db string) []string {
	rg.locker.Lock()
	defer rg.locker.Unlock()
	response := getLoadDataResponse(rg, db)
	rg.data = response.SQLs
	rg.query = response.Queries
	return rg.data
}

func (rg *RandGen) LoadQuery(db string) []string {
	rg.locker.Lock()
	defer rg.locker.Unlock()
	return rg.query
}

func (rg *RandGen) Name() string {
	return "randgen"
}

func (rg *RandGen) GetDataLoaders() []interfaces.DataLoader {
	return []interfaces.DataLoader{rg}
}

func (rg *RandGen) GetQueryLoaders() []interfaces.QueryLoader {
	return []interfaces.QueryLoader{rg}
}

func (rg *RandGen) GetComparor() interfaces.SqlResultComparer {
	return &sqldiff.StandardComparer{}
}

func (rg *RandGen) GetLoop() int {
	return rg.Loop
}

func (rg *RandGen) GetLoopInterval() int {
	return rg.LoopInterval
}

func getLoadDataResponse(randgen *RandGen, db string) *LoadDataResponse {
	payload := &LoadDataRequest{Yy: randgen.Yy, ZZ: randgen.ZZ, DB: db, Queries: randgen.Queries}
	resp, err := http.Post("http://localhost:9080/loaddata", "application/json",
		strings.NewReader(getLoadDataRequestString(payload)))
	if err != nil {
		fmt.Println(err)
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
