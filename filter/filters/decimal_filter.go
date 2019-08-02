package filter

import (
	"database/sql"
	"github.com/dqinyuan/sqlparser"
	"strings"
)

type DecimalFilter struct {
}

func (f *DecimalFilter) CompareHook(cv1 string, cv2 string, colType *sql.ColumnType) bool {
	return false
}

func (f *DecimalFilter) ErrType() string {
	return "default decimal precision not consistent"
}

func (f *DecimalFilter) Init() {
}

func (f *DecimalFilter) ClearFilterCtx() {
}

func (f *DecimalFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {
	switch statement := statement.(type) {
	case *sqlparser.DDL:
		if statement.TableSpec == nil || statement.TableSpec.Columns == nil {
			return statement
		}

		// 遍历所有字段声明
		for _, cd := range statement.TableSpec.Columns {
			if strings.ToLower(cd.Type.Type) == "decimal" && cd.Type.Length == nil {
				cd.Type.Length = &sqlparser.SQLVal{
					Type: sqlparser.IntVal,
					Val:  []byte("10"),
				}
			}
		}
	}

	return statement
}

func (f *DecimalFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
}





