package filter

import (
	"database/sql"
	"encoding/json"
	"github.com/dqinyuan/sqlparser"
	"reflect"
	"strings"
)

// filter json col
type JsonFilter struct {
}

func (*JsonFilter) Init() {
}

func (*JsonFilter) FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement) {
	return statement
}

func (*JsonFilter) CompareHook(cv1 string, cv2 string, colType *sql.ColumnType) bool {
	//maybe it's json
	if (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) ||
		(strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) {
		if jsonEquals(cv1, cv2) {
			return true
		}
	}
	return false
}

func jsonEquals(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func (*JsonFilter) UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
}

func (*JsonFilter) ClearFilterCtx() {
}

func (*JsonFilter) ErrType() string {
	return "json col"
}


