package filter

import (
	"database/sql"
	"fmt"
	"github.com/monitoring-system/dbtest/config"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	_ "time"
)

const (
	FilterPATH = "plugin-filters"
)

type ErrMsgFilter func(errMsg string, source string) bool

type DiffFilter func(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool

var (
	errMsgFilters = make([]ErrMsgFilter, 0)
	diffFilters   = make([]DiffFilter, 0)
	mux           sync.Mutex
)

func RegisterErrMsgFilter(f ErrMsgFilter) {
	mux.Lock()
	defer mux.Unlock()

	errMsgFilters = append(errMsgFilters, f)
}

func RegisterDiffFilter(f DiffFilter) {
	mux.Lock()
	defer mux.Unlock()

	diffFilters = append(diffFilters, f)
}

func FilterError(errMsg string, source string) bool {
	if config.Conf.TraceAllErrors {
		return false
	}
	for _, f := range errMsgFilters {
		if f(errMsg, source) {
			return true
		}
	}

	return false
}

func FilterCompareDiff(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool {
	for _, f := range diffFilters {
		if f(vInTiDB, vInMySQL, colType) {
			return true
		}
	}

	return false
}

var (
	errMsgFilter ErrMsgFilter
	diffFilter   DiffFilter

	errMsgType reflect.Type
	diffType   reflect.Type
)

func init() {
	os.Mkdir(FilterPATH, os.ModePerm)

	errMsgType = reflect.TypeOf(errMsgFilter)
	diffType = reflect.TypeOf(diffFilter)

	// Load filters from FilterPATH.
	if files, err := ioutil.ReadDir(FilterPATH); err != nil {
		log.Error("load filters failed", zap.Error(err))
	} else {
		for _, info := range files {
			if err := AddFilter(fmt.Sprintf("%s/%s", FilterPATH, info.Name())); err != nil {
				log.Error("load filter failed", zap.Error(err))
			}
		}
	}

	// Load default filters
	loadDefaultErrMsgFilters()
	loadDefaultDiffFilters()
}

func loadDefaultErrMsgFilters() {

}

func loadDefaultDiffFilters() {
	RegisterDiffFilter(filterNumberPercision)
	RegisterDiffFilter(filterZero)
}

func filterNumberPercision(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool {
	if !isFloat(colType) {
		return false
	}
	tidb, err := floatConvert(vInTiDB.([]byte))
	if err != nil {
		return false
	}

	mysql, err := floatConvert(vInMySQL.([]byte))
	if err != nil {
		return false
	}

	if floatRound(tidb.(float64), 4) == floatRound(mysql.(float64), 4) {
		return true
	}

	return false
}

func floatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func filterZero(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool {
	f, ok := typeForMysqlToGo[strings.ToLower(colType.DatabaseTypeName())]
	if !ok {
		log.Warn("Unsupport type", zap.String("type", colType.DatabaseTypeName()))
		return false
	}

	tidb, err := f(vInTiDB.([]byte))
	if err != nil {
		log.Warn("Convert function failed", zap.Error(err))
		return false
	}

	mysql, err := f(vInMySQL.([]byte))
	if err != nil {
		log.Warn("Convert function failed", zap.Error(err))
		return false
	}

	if tidb == mysql {
		return true
	}

	if isZero(tidb) && isZero(mysql) {
		return true
	}

	return false
}

func isZero(v interface{}) bool {

	if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
		return true
	}

	return false
}

type convertF func(data []byte) (interface{}, error)

func intConvert(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 64)
	return v, err
}

func boolConvert(data []byte) (interface{}, error) {
	v, err := strconv.ParseBool(string(data))
	return v, err
}

func floatConvert(data []byte) (interface{}, error) {
	v, err := strconv.ParseFloat(string(data), 64)
	return v, err
}

func stringConvert(data []byte) (interface{}, error) {
	return string(data), nil
}

//func timeConvert(data []byte) (interface{}, error) {
//
//	time.
//}

var typeForMysqlToGo = map[string]convertF{
	"int":                intConvert,
	"integer":            intConvert,
	"tinyint":            intConvert,
	"smallint":           intConvert,
	"mediumint":          intConvert,
	"bigint":             intConvert,
	"int unsigned":       intConvert,
	"integer unsigned":   intConvert,
	"tinyint unsigned":   intConvert,
	"smallint unsigned":  intConvert,
	"mediumint unsigned": intConvert,
	"bigint unsigned":    intConvert,
	"bit":                intConvert,
	"bool":               boolConvert,
	"enum":               stringConvert,
	"set":                stringConvert,
	"varchar":            stringConvert,
	"char":               stringConvert,
	"tinytext":           stringConvert,
	"mediumtext":         stringConvert,
	"text":               stringConvert,
	"longtext":           stringConvert,
	"blob":               stringConvert,
	"tinyblob":           stringConvert,
	"mediumblob":         stringConvert,
	"longblob":           stringConvert,
	//"date":               reflect.TypeOf(time.Now()), // time.Time or string
	//"datetime":           reflect.TypeOf(time.Now()), // time.Time or string
	//"timestamp":          reflect.TypeOf(time.Now()), // time.Time or string
	//"time":               reflect.TypeOf(time.Now()), // time.Time or string
	"float":     floatConvert,
	"double":    floatConvert,
	"decimal":   floatConvert,
	"binary":    stringConvert,
	"varbinary": stringConvert,
}

func isFloat(colType *sql.ColumnType) bool {
	t := strings.ToLower(colType.DatabaseTypeName())
	return t == "float" || t == "double" || t == "decimal"
}
