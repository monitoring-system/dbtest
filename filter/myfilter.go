package filter

import (
	"github.com/xwb1989/sqlparser"
	"sync"
)

type Filter interface {
	Init()
	// filter by sql run time context, should return nil if sql is filtered
	FiltByCtx(statement sqlparser.Statement) (newStmt sqlparser.Statement)
	// update context by sql exec result
	UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error)
	// clear filter context
	ClearFilterCtx()
	// err type description
	ErrType() string
}

var filtersLock = &sync.RWMutex{}
var filters = make([]Filter, 0)

func Register(filter Filter) {
	filtersLock.Lock()
	defer filtersLock.Unlock()
	filter.Init()
	filters = append(filters, filter)
}

// Sequencial Filter
// Perhaps problem: sql rewrite conflict
func traverse(f func(filter Filter) bool) {
	filtersLock.RLock()
	defer filtersLock.RUnlock()

	for _, filter := range filters {
		if !f(filter) {
			break
		}
	}
}

func FiltByCtx(statement sqlparser.Statement) (smt sqlparser.Statement, comment string) {
	traverse(func(filter Filter) bool {
		statement = filter.FiltByCtx(statement)
		if statement == nil {
			// break traverse
			comment = filter.ErrType()
			return false
		}
		return true
	})
	return statement, comment
}

func UpdateCtxByExecResult(statement sqlparser.Statement, err0 error, err1 error) {
	traverse(func(filter Filter) bool {
		filter.UpdateCtxByExecResult(statement, err0, err1)
		return true
	})
}

func ClearFilterCtx() {
	traverse(func(filter Filter) bool {
		filter.ClearFilterCtx()
		return true
	})
}

func init() {
	/*	os.Mkdir(FilterPATH, os.ModePerm)

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
		loadDefaultDiffFilters()*/
		loadDefaultFilters()
}

func loadDefaultFilters() {
	Register(&DecimalFilter{})
	Register(&FloatFilter{})
	Register(&TableNotExistFilter{})
}