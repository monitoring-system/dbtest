package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestPostTests(t *testing.T) {
	t.SkipNow()
	config.Conf.StandardDB = "root:123456@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local"

	go StartServer()

	time.Sleep(5 * time.Second)

	testSubmit := &types.Test{
		SqlLoader: "string",
		DataStr: "drop table if exists nums;create table nums (num float);insert into nums (num) values(10.123456);" +
			"select * from nums",
		DbRetain: true,
	}

	jBytes, err := json.Marshal(testSubmit)
	if err != nil {
		t.Fatalf("can not marshal to json %v\n", err)
	}

	resp, err := http.Post("http://127.0.0.1:8080/tests", "application/json", bytes.NewBuffer(jBytes))
	if err != nil {
		t.Fatalf("reponse error %v\n", err)
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read error %v\n", err)
	}

	fmt.Println(string(all))

	time.Sleep(1000 * time.Second)
}
