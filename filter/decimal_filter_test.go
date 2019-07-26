package filter

import (
	"github.com/stretchr/testify/assert"
	"github.com/xwb1989/sqlparser"
	"testing"
)

func TestDecimalTest(t *testing.T) {
	d := &DecimalFilter{}

	sql := "create table dt (t decimal);"

	ast, _ := sqlparser.Parse(sql)

	newAst := d.FiltByCtx(ast)
	expected := `create table dt (
	t decimal(10)
)`
	assert.Equal(t, sqlparser.String(newAst), expected)
}
