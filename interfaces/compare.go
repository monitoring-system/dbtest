package interfaces

import (
	"database/sql"
)

type SqlResultComparer interface {
	// execute query on two databases and compare the result
	CompareQuery(db1, db2 *sql.DB, query string) (bool, error, error)
}

type CellFilter interface {
	// return true if we can ignore the difference of two cells
	Filter([]byte, []byte, *sql.ColumnType) bool
}
