package filter

import (
	"github.com/pkg/errors"
	"github.com/dqinyuan/sqlparser"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

// sql parser not support 'create ... as ...' and 'create ... like ...'

func TestTableNotExistFilter(t *testing.T) {
	f := &TableNotExistFilter{}

	f.Init()
	defer f.ClearFilterCtx()

	stmt := runSql(t, f, "select t from f;", errors.New(""), nil)
	assert.NotEqual(t, stmt, nil)

	stmt = runSql(t, f, "select t from f;", nil, nil)
	assert.NotEqual(t, stmt, nil)

	stmt = runSql(t, f, "create table t (f float);", nil, nil)
	assert.NotEqual(t, stmt, nil)

	stmt = runSql(t, f, "select t from f;", nil, nil)
	assert.NotEqual(t, stmt, nil)

	stmt = runSql(t, f, "create table t (f float);", errors.New(""), nil)
	assert.NotEqual(t, stmt, nil)

	stmt = runSql(t, f, "select t from f;", nil, nil)
	assert.Equal(t, stmt, nil)
}

func runSql(t *testing.T, f Filter, sql string, err1 error, err2 error) sqlparser.Statement {
	statement, err := sqlparser.Parse(sql)
	if err != nil {
		t.Fatalf("sql parse error %v\n", err)
	}

	newStmt := f.FiltByCtx(statement)
	if newStmt == nil {
		return newStmt
	}

	f.UpdateCtxByExecResult(newStmt, err1, err2)

	return newStmt
}



