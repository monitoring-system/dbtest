package filter

import (
	"github.com/prometheus/common/log"
	"k8s.io/client-go/util/retry"
	"sync"
)

type Filter interface {

}

type FilterErrMsg interface {
	Ignore(errMsg string, source string) bool
}

type FilterDiff interface {
	Ignore(vInTiDB interface{}, vInMySQL interface{}) bool
}


type defaultFilterErrMsg struct {

}

func (f defaultFilterErrMsg) Ignore(errMsg string, source string) bool {
	code, msg := decodeMsg(errMsg)
	return getFilterAndInsertIfNotExist(code, msg, source)
}

func decodeMsg(errMsg string) (int, string)  {
	return -1, ""
	
}

type defaultFilterDiff struct {

}

func (f defaultFilterDiff) Ignore(vInTiDB interface{}, vInMySQL interface{}) bool {
	if v, ok := vInTiDB.(string); ok {

	}
}

var (
	errMsgFilters = make([]FilterErrMsg, 0)
	diffFilters = make([]FilterDiff, 0)
	mux sync.Mutex
)

func RegisterFilters() {
	mux.Lock()
	defer mux.Unlock()

	errMsgFilters = []FilterErrMsg{defaultFilterErrMsg{}}
	diffFilters = []FilterDiff{defaultFilterDiff{}}





}

func RegisterFilter(f Filter) {
	if v, ok := f.(FilterErrMsg); ok {
		errMsgFilters = append(errMsgFilters, v)
		return
	}

	if v, ok := f.(FilterDiff); ok {
		diffFilters = append(diffFilters, v)
		return
	}

	log.Warn("Unknown filter type")
}

func UnResiterFilters() {
	mux.Lock()
	defer mux.Unlock()

	errMsgFilters = []FilterErrMsg{defaultFilterErrMsg{}}
	diffFilters = []FilterDiff{defaultFilterDiff{}}
}

func FilterError(errMsg string, source string) bool{
	for _, f := range errMsgFilters {
		if f.Ignore(errMsg, source) {
			return true
		}
	}

	return false
}

func FilterCompareDiff(vInTiDB interface{}, vInMySQL interface{}) bool {
	for _, f := range diffFilters {
		if f.Ignore(vInTiDB, vInMySQL) {
			return true
		}
	}

	return false
}