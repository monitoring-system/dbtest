package impl

import (
	"github.com/monitoring-system/dbtest/randgen"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"sync"
)

const LocalRandgen = "local-randgen"

type LocalRandgenLoader struct {
	Yy      string
	Zz      string
	Queries int

	locker sync.RWMutex
	data   []string
	query  []string
}

func (test *LocalRandgenLoader) LoadData(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	test.getLoadDataResponse(db)
	return test.data
}

func (test *LocalRandgenLoader) LoadQuery(db string) []string {
	test.locker.Lock()
	defer test.locker.Unlock()
	return test.query
}

func (test *LocalRandgenLoader) Name() string {
	return LocalRandgen
}

func (test *LocalRandgenLoader) getLoadDataResponse(db string) {
	loader := &randgen.SQLGenerator{}
	loader.Init(db)
	sqls, err := loader.LoadData(test.Zz, test.Yy, db, test.Queries)
	if err != nil {
		log.Info("generate data failed", zap.Error(err))
	}
	test.data = sqls
	test.query = loader.CachedQueries
}
