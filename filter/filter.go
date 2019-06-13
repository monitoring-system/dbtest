package filter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
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
	diffFilter	DiffFilter

	errMsgType reflect.Type
	diffType reflect.Type
)

func init() {
	os.Mkdir(FilterPATH, os.ModePerm)

	errMsgType = reflect.TypeOf(errMsgFilter)
	diffType = reflect.TypeOf(diffFilter)

	// Load filters from FilterPATH.
	if files, err := ioutil.ReadDir(FilterPATH); err != nil {
		log.Error("load filters failed", zap.Error(err))
	}else {
		for _, info := range files {
			if err := AddFilter(fmt.Sprintf("%s/%s",FilterPATH, info.Name())); err != nil {
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
	if strings.ToUpper(colType.DatabaseTypeName()) != "DECIMAL" {
		return false
	}

	if floatRound(toFload64(vInTiDB), 4) == floatRound(toFload64(vInMySQL), 4) {
		return true
	}

	return false
}

func floatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

func toFload64(v interface{}) float64{
	f, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		log.Error("convert float64 failed", zap.Error(err))
		return 0
	}

	return f
}

func filterZero(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool {
	t, ok := typeForMysqlToGo[strings.ToLower(colType.DatabaseTypeName())]
	if !ok {
		log.Warn("Unsupport type", zap.String("type", colType.DatabaseTypeName()))
		return false
	}


	if isZero(vInTiDB.([]byte), t) && isZero(vInMySQL.([]byte), t) {
		return true
	}

	 //colType.ScanType()
	 return false
}

func isZero(v []byte, t reflect.Type) bool {
	value := reflect.Zero(t).Interface()
	json.Unmarshal(v, &value)
	if value == reflect.Zero(t).Interface() || value == nil {
		return true
	}

	return false
}

var typeForMysqlToGo = map[string]reflect.Type{
	"int":                reflect.TypeOf(0),
	"integer":            reflect.TypeOf(0),
	"tinyint":            reflect.TypeOf(0),
	"smallint":           reflect.TypeOf(0),
	"mediumint":          reflect.TypeOf(0),
	"bigint":             reflect.TypeOf(0),
	"int unsigned":       reflect.TypeOf(0),
	"integer unsigned":   reflect.TypeOf(0),
	"tinyint unsigned":   reflect.TypeOf(0),
	"smallint unsigned":  reflect.TypeOf(0),
	"mediumint unsigned": reflect.TypeOf(0),
	"bigint unsigned":    reflect.TypeOf(0),
	"bit":                reflect.TypeOf(0),
	"bool":               reflect.TypeOf(false),
	"enum":               reflect.TypeOf(""),
	"set":                reflect.TypeOf(""),
	"varchar":            reflect.TypeOf(""),
	"char":               reflect.TypeOf(""),
	"tinytext":           reflect.TypeOf(""),
	"mediumtext":         reflect.TypeOf(""),
	"text":               reflect.TypeOf(""),
	"longtext":           reflect.TypeOf(""),
	"blob":               reflect.TypeOf(""),
	"tinyblob":           reflect.TypeOf(""),
	"mediumblob":         reflect.TypeOf(""),
	"longblob":           reflect.TypeOf(""),
	"date":               reflect.TypeOf(time.Now()), // time.Time or string
	"datetime":           reflect.TypeOf(time.Now()), // time.Time or string
	"timestamp":          reflect.TypeOf(time.Now()), // time.Time or string
	"time":               reflect.TypeOf(time.Now()), // time.Time or string
	"float":              reflect.TypeOf(float64(0)),
	"double":             reflect.TypeOf(float64(0)),
	"decimal":            reflect.TypeOf(float64(0)),
	"binary":             reflect.TypeOf(""),
	"varbinary":          reflect.TypeOf(""),
}



