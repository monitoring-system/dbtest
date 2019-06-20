package main

import (
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/randgen"
	"github.com/satori/go.uuid"
	"net/http"
)

func main() {
	StartServer()
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/loaddata", LoadData)
	err := http.ListenAndServe(":9080", mux)
	if err != nil {
		fmt.Println(err)
	}
}

func MustJosnMarshal(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	payload := LoadDataRequest{}
	json.NewDecoder(r.Body).Decode(&payload)

	loader := &randgen.SQLGenerator{}
	loader.Init(uuid.NewV4().String())
	fmt.Println(payload.ZZ)
	fmt.Println(payload.Yy)

	sqls, err := loader.LoadData(payload.ZZ, payload.Yy, payload.DB, payload.Queries)
	if err != nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
	}
	w.Write(MustJosnMarshal(&LoadDataResponse{SQLs: sqls, Queries: loader.CachedQueries}))
}

type LoadDataRequest struct {
	Yy      string `json:"yy"`
	ZZ      string `json:"zz"`
	DB      string `json:"db"`
	Queries int    `json:"queries"`
}

type LoadDataResponse struct {
	SQLs    []string `json:"sql"`
	Queries []string `json:"queries"`
}

type RandgenLoader struct {
	TestName string
	// 存放yyzz文件的path
	ConfPath string
	//  存放结构的path
	ResultPath string
	//  randgen主目录
	RmPath string

	CachedQueries []string
}
