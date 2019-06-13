package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	tables, err := Parse("select * from t1 as t left join t2 tt on t.a = tt.a")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, t := range  tables.TableName {
		fmt.Println(t)
	}
}
