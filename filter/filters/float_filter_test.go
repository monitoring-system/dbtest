package filter

import (
	"fmt"
	"github.com/dqinyuan/sqlparser"
	"testing"
)

func TestFloatFilter(t *testing.T) {

	f := &FloatFilter{}

	sql := "create table floattest(f float);"
	statement, _ := sqlparser.Parse(sql)

	newAst := f.FiltByCtx(statement)

	fmt.Println(sqlparser.String(newAst))
}
