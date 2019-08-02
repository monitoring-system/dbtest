package filter

import (
	"github.com/dqinyuan/sqlparser"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestLogicFloatFilter(t *testing.T) {
	sqls := []string{
		"select null || 0.1;",
		"select null && 0.1;",
		"select null or 0.1;",
		"select null and 0.1;",
		"select col_double_key and null;",
		"select col_float_key or null;"}
	f := &LogicFloatFilter{}

	for _, s := range sqls {
		ast, _ := sqlparser.Parse(s)
		newStmt := f.FiltByCtx(ast)
		assert.Equal(t, newStmt, nil)
	}
}
