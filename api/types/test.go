package types

import (
	"encoding/json"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/monitoring-system/dbtest/sqldiff"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Test struct {
	ID           int64  `json:"id",gorm:"primary_key" rql:"filter,sort"`
	TestName     string `json:"testName" rql:"filter,sort"`
	DataLoader   string `json:"dataLoader" rql:"filter,sort"`
	QueryLoader  string `json:"queryLoader" rql:"filter,sort"`
	Loop         int    `json:"loop"`
	LoopInterval int    `json:"loopInterval" rql:"filter,sort"`

	Yy      string `gorm:"type:TEXT;" rql:"filter,sort"`
	Zz      string `gorm:"type:TEXT;"`
	Queries int    `json:"queries"`

	locker sync.RWMutex
	data   []string
	query  []string
}

//persistent the result and set the id
func (test *Test) Persistent() error {
	return db.GetDB().Create(test).Error
}

func (test *Test) Update() error {
	return db.GetDB().Save(test).Error
}

func GetTestById(id int64) (*Test, error) {
	result := &Test{}
	return result, db.GetDB().Where(" id=?", id).First(result).Error
}

func ListTest() ([]*Test, error) {
	var list []*Test
	return list, db.GetDB().Find(&list).Error
}

func (test *Test) LoadData(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	response := getLoadDataResponse(test, db)
	test.data = response.SQLs
	test.query = response.Queries
	return test.data
}

func (test *Test) LoadQuery(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	return test.query
}

func (test *Test) Name() string {
	return "randgen"
}

func (test *Test) GetName() string {
	return ""
}

func (test *Test) GetDataLoaders() []interfaces.DataLoader {
	return []interfaces.DataLoader{test}
}

func (test *Test) GetQueryLoaders() []interfaces.QueryLoader {
	return []interfaces.QueryLoader{test}
}

func (test *Test) GetComparor() interfaces.SqlResultComparer {
	return &sqldiff.StandardComparer{}
}

func (test *Test) GetLoop() int {
	return test.Loop
}

func (test *Test) GetLoopInterval() int {
	return test.LoopInterval
}

func getLoadDataResponse(randgen *Test, db string) *LoadDataResponse {
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
