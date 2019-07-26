package sqldiff

import (
	"database/sql"
	"fmt"
	"testing"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetQueryResult(t *testing.T) {
	t.SkipNow()
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("open db error %v\n", err)
	}

	sqlResult, err := GetQueryResult(db, "insert into t values (10);")
	if err != nil {
		t.Fatalf("statement error %v\n", err)
	}
	fmt.Println(sqlResult.header) // []
	fmt.Println(sqlResult.data) // []
}
