package util

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/log"
	"time"
)

// OpenDBWithRetry opens a database specified by its database driver name and a
// driver-specific data source name. And it will do some retries if the connection fails.
func OpenDBWithRetry(driverName, dataSourceName string) (mdb *sql.DB, err error) {
	startTime := time.Now()
	sleepTime := time.Millisecond * 500
	retryCnt := 60
	// The max retry interval is 30 s.
	for i := 0; i < retryCnt; i++ {
		mdb, err = sql.Open(driverName, dataSourceName)
		if err != nil {
			fmt.Printf("open db failed, retry count %d err %v\n", i, err)
			time.Sleep(sleepTime)
			continue
		}
		err = mdb.Ping()
		if err == nil {
			break
		}
		log.Warnf("ping db failed, retry count %d err %v", i, err)
		mdb.Close()
		time.Sleep(sleepTime)
	}
	if err != nil {
		log.Errorf("open db failed %v, take time %v", err, time.Since(startTime))
		return nil, err
	}

	return
}
