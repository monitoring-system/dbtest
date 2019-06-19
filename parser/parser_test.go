package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	sql := "create table t(id decimal)"
	tables, err := Parse(sql)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if tables.Rewrite {
		sql = tables.NewSql
	}
	fmt.Println(sql)
	for _, t := range tables.TableName {
		fmt.Println(t)
	}
}
