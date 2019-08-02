package filter

import (
	"database/sql"
	"github.com/dqinyuan/sqlparser"
	"github.com/pkg/errors"
	"strings"
)

/*
过滤float类型参与逻辑(and or)运算的Bug
比如 select null && 0.1;
*/

type LogicFloatFilter struct {
}

func (*LogicFloatFilter) Init() {
}

func (*LogicFloatFilter) checkColName(e sqlparser.Expr) bool  {
	colName, ok := e.(*sqlparser.ColName)
	if !ok {
		return false
	}

	return strings.Contains(colName.Name.String(),
		"double") || strings.Contains(colName.Name.String(), "float")
}

// false表示非浮点数, true表示是浮点数
func (f *LogicFloatFilter) checkSQLVal(e sqlparser.Expr) bool  {
	// val will be nil when type assertion fail
	val, ok := e.(*sqlparser.SQLVal)
	if !ok {
		return f.checkColName(e)
	}

	if val.Type == sqlparser.FloatVal {
		return true
	}

	return false
}

func (f *LogicFloatFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {
	ignore := false
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.AndExpr:
			ignore = f.checkSQLVal(n.Left) || f.checkSQLVal(n.Right)
		case *sqlparser.OrExpr:
			ignore = f.checkSQLVal(n.Left) || f.checkSQLVal(n.Right)
		}

		if ignore {
			return false, errors.New("")
		}

		return true, nil
	}, statement)

	if ignore {
		return nil
	} else  {
		return statement
	}
}

func (*LogicFloatFilter) CompareHook(cv1 string, cv2 string, colType *sql.ColumnType) bool {
	return false
}

func (*LogicFloatFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
}

func (*LogicFloatFilter) ClearFilterCtx() {
}

func (*LogicFloatFilter) ErrType() string {
	return "float part in logic computation"
}



