package executor

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/monitoring-system/dbtest/api/types"
	_ "github.com/monitoring-system/dbtest/filter/filters"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/stretchr/testify/assert"
	golog "log"
	"os"
	"strings"
	"testing"
)

func TestNumberPrecisionFilter(t *testing.T) {

	runSqls := []string{
		"create table nums (num int, f float)",
		"create table nums (num int, d decimal)",
	}
	expects := []string {
		"create table nums (num int,f float(7,4))",
		"create table nums (num int,d decimal(10))",
	}

	mockCtl := gomock.NewController(t)
	mockCompare := interfaces.NewMockSqlResultComparer(mockCtl)
	// len(sqls) + 1  because stringLoader will add 'use db' atomically
	mockCompare.EXPECT().CompareQuery(nil, nil,
		gomock.Any()).Return("", nil, nil).AnyTimes()

	testFilter(t, runSqls, expects, mockCompare)
}

func TestTableNotExistsFilter(t *testing.T) {
	runSqls := []string{
		"insert into nums values(1, 4)",
		"create table nums (num int, f int)",
		"insert into nums values(1, 4)",
		"insert into tels values(1, 4)",
	}
	expects := []string{
		"insert into nums values (1, 4)",
		"create table nums (num int,f int)",
		"",   //filted
		"insert into tels values (1, 4)",
	}


	mockCtl := gomock.NewController(t)
	mockCompare := interfaces.NewMockSqlResultComparer(mockCtl)
	// len(sqls) + 1  because stringLoader will add 'use db' atomically
	mockCompare.EXPECT().CompareQuery(nil, nil,
		gomock.Eq("create table nums (num int,f int)")).Return("", errors.New("create err"), nil)
	mockCompare.EXPECT().CompareQuery(nil, nil,
		gomock.Any()).Return("", nil, nil).AnyTimes()

	testFilter(t, runSqls, expects, mockCompare)
}

// expectedSqls中的空字符串表示期待其被过滤
func testFilter(t *testing.T, inputSqls []string, expectedSqls []string, mockCompare interfaces.SqlResultComparer) {

	filteredCounter := 0
	moves := make([]int, len(inputSqls))
	for i := range moves {
		if expectedSqls[i] == "" {
			filteredCounter++
		}
		moves[i] = filteredCounter
	}

	test := &types.Test{
		SqlLoader:"string",
		DataStr:strings.Join(inputSqls, ";"),
	}

	scope := &testScope{test:test, loopResult:&types.LoopResult{}, logger: golog.New(os.Stdout,
		"", golog.LstdFlags), dbName:"test", db1: nil, db2: nil}

	exec := &Executor{}

	eSqls := exec.execSql(scope, mockCompare)
	sqlArr := strings.Split(eSqls.String(), ";\n")[1:]

	for i := range inputSqls {
		if expectedSqls[i] == "" {
			continue
		}

		assert.Equal(t, expectedSqls[i], sqlArr[i - moves[i]])
	}
}