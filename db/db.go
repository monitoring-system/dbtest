package db

import (
	"fmt"
	"github.com/prometheus/common/log"

	"github.com/jinzhu/gorm"
)

const (
	filterTableName = "filters"
)

var db *gorm.DB

func InitDatabase(dsn string, models []interface{}) {
	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("init db failed, err=%ver", err))
	}
	for _, model := range models {
		db.AutoMigrate(model)
	}
}

func GetDB() *gorm.DB {
	return db
}

func GetFilterAndInsertIfNotExist(errcode int, keyword string, source string) bool {
	r, err := db.DB().Query(fmt.Sprintf("select id from %s where filter=?", filterTableName), buildFilterField(errcode, keyword))
	if err != nil {
		log.Warn("query filter failed", err)
		return false
	}

	var id int
	for r.Next() {
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
	_, err := db.DB().Exec(fmt.Sprintf("insert into %s (filter, source) values(?, ?)", filterTableName), buildFilterField(errcode, keyword), sql)
	if err != nil {
		log.Warn("insert filter failed", err)
	}
}

func buildFilterField(errcode int, keyword string) string {
	return fmt.Sprintf("%d_%s", errcode, keyword)
}
