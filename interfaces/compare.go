package interfaces

import (
	"database/sql"
)

type SqlResultComparer interface {
	// execute query on two databases and compare the result
	CompareQuery(db1, db2 *sql.DB, query string) (string, error, error)
}
