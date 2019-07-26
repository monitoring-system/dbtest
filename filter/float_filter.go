package filter

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

type FloatFilter struct {
}

func (f *FloatFilter) ErrType() string {
	return "default float precision not consistent"
}

func (f *FloatFilter) Init() {
}

func (f *FloatFilter) ClearFilterCtx() {
}

func (f *FloatFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {
	switch statement := statement.(type) {
	case *sqlparser.DDL:
		if statement.TableSpec == nil || statement.TableSpec.Columns == nil {
			return statement
		}

		// 遍历所有字段声明
		for _, cd := range statement.TableSpec.Columns {
			if strings.ToLower(cd.Type.Type) == "float" && cd.Type.Length == nil {
				cd.Type.Length = &sqlparser.SQLVal{
					Type: sqlparser.IntVal,
					Val:  []byte("7"),
				}

				cd.Type.Scale = &sqlparser.SQLVal{
					Type: sqlparser.IntVal,
					Val:  []byte("4"),
				}
			}
		}
	}

	return statement
}

func (f *FloatFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
}

