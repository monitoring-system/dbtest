package filter

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/common/log"
)

const (
	datasourceName  = "test.db"
	tableName = "filters"
)

var (
	db *sql.DB
	createTableSQL = fmt.Sprintf("create if not exists table %s (id int primary key AUTOINCREMENT, filter varchar(256) unique key, source text)", tableName)
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", datasourceName)
	if err != nil {
		panic("init db failed")
	}

	db.Exec(createTableSQL)
}

func getFilterAndInsertIfNotExist(errcode int, keyword string, source string) bool{
	r, err := db.Query(fmt.Sprintf("select id from %s where filter=?", tableName), buildFilterField(errcode, keyword))
	if err != nil {
		log.Warn("query filter failed", err)
		return false
	}

	var id int
	for r.Next()  {
		err := r.Scan(&id)
		if err != nil {
			log.Warn("query filter failed", err)
		}
		break
	}

	if id > 0 {
		return true
	}

	insertFilter(errcode, keyword, source)
	return false
}

func insertFilter(errcode int, keyword string, sql string) {
	_, err := db.Exec(fmt.Sprintf("insert into %s values (filter, source) values(?, ?)", tableName), buildFilterField(errcode, keyword), sql)
	if err != nil {
		log.Warn("insert filter failed", err)
	}
}

func buildFilterField(errcode int, keyword string) string{
	return fmt.Sprintf("%d_%s", errcode, keyword)
}
