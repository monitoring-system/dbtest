package sqldiff

import (
	"database/sql"
)

type StandardComparer struct {
	Strict bool
}

func (c *StandardComparer) CompareQuery(db1, db2 *sql.DB, query string) (bool, error, error) {
	result1, err1 := GetQueryResult(db1, query)
	result2, err2 := GetQueryResult(db2, query)
	if err1 != nil || err2 != nil {
		return false, err1, err2
	} else {
		// now compare the results
		equals := false
		if c.Strict {
			equals = result2.strictCompare(result1)
		} else {
			equals = result2.nonOrderCompare(result1)
		}
		return equals, nil, nil
	}
}
