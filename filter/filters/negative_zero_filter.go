package filter

import (
	"database/sql"
	"github.com/dqinyuan/sqlparser"
)

// filter -0 bug in Mysql
type NegZeroFilter struct {
}

func (*NegZeroFilter) Init() {
}

func (*NegZeroFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {
	return statement
}

func (*NegZeroFilter) CompareHook(cv1 string, cv2 string, colType *sql.ColumnType) bool {
	if cv1 == "-0" || cv2 == "-0" {
		return true
	}
	return false
}

func (*NegZeroFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
}

func (*NegZeroFilter) ClearFilterCtx() {
}

func (*NegZeroFilter) ErrType() string {
	return "mysql bug -0"
}

