package impl

import (
	"database/sql"
	"github.com/monitoring-system/dbtest/sqldiff"
	"strings"
)

type ContainsComparer struct {
	Content string
}

//execute on tidb and get all rows data as a single string then check if that string contains the specified Content
//can be used to check the tidb explain result
func (c *ContainsComparer) CompareQuery(db1, db2 *sql.DB, query string) (string, error, error) {
	result, err := sqldiff.GetQueryResult(db2, query)
	if err != nil {
		return "", nil, err
	}
	resultStr := result.String()
	if strings.Contains(resultStr, c.Content) {
		return "", nil, nil
	} else {
		return sqldiff.GetColorDiff(c.Content, resultStr), nil, nil
	}
}
