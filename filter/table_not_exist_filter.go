package filter

import (
	"github.com/monitoring-system/dbtest/util"
	"github.com/pkg/errors"
	"github.com/xwb1989/sqlparser"
)

/*
if ddl error, we will ignore all sql subsequently about the table
 */
type TableNotExistFilter struct {
	ignoreTables *util.Set
}

func (f *TableNotExistFilter) ErrType() string {
	return "previous relevant ddl not run successfully in both db"
}

func (f *TableNotExistFilter) Init() {
	f.ignoreTables = util.NewSet()
}

func (f *TableNotExistFilter) ClearFilterCtx() {
	f.ignoreTables = util.NewSet()
}

// if need ignore by ignore tables
func (f *TableNotExistFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {

	ignore := false

	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.DDL:
			return false, nil
		case sqlparser.TableName:
			if f.ignoreTables.Contains(n.Name.String()) {
				ignore = true
				return false, errors.New("")
			}
		}
		return true, nil
	}, statement)

	if ignore {
		return nil
	} else {
		return statement
	}
}

// update ignore tables
func (f *TableNotExistFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
	if (err0 != nil && err1 == nil) || (err0 == nil && err1 != nil) {

		_, ok := statement.(*sqlparser.DDL)

		if ok {
			sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				switch n := node.(type) {
				case sqlparser.TableName:
					f.ignoreTables.Put(n.Name.String())
				}
				return true, nil
			}, statement)
		}
	}
}

