package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	t.Run("test decimal rewrite", func(t *testing.T) {
		sql := "create table t(id decimal)"
		tables, err := Parse(sql)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if tables.Rewrite {
			sql = tables.NewSql
		} else {
			fmt.Println("should rewrite")
			t.FailNow()
		}

		expectedSql :=
			`create table t (
	id decimal(10)
)`
		if sql != expectedSql {
			t.Errorf("%s\n not expected\n", sql)
		}

		for _, tt := range tables.TableName {
			if tt != "t" {
				t.Errorf("%s table name not expected", tt)
			}
		}
	})

	t.Run("test insert ignore rewrite", func(t *testing.T) {
		sql := "INSERT /*! IGNORE */ INTO t VALUES (10)"
		tables, err := Parse(sql)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}

		if tables.Rewrite {
			sql = tables.NewSql
		} else {
			fmt.Println("should rewrite")
			t.FailNow()
		}

		expectedSql := `insert into t values (10)`

		if sql != expectedSql {
			t.Errorf("%s\n not expected\n", sql)
		}

	})

}
