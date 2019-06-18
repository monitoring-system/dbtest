package types

import (
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/interfaces"
	"github.com/monitoring-system/dbtest/interfaces/impl"
	"github.com/monitoring-system/dbtest/sqldiff"
	"sync"
)

type Test struct {
	ID           int64  `json:"id",gorm:"primary_key" rql:"filter,sort"`
	TestName     string `json:"testName" rql:"filter,sort"`
	DataLoader   string `json:"dataLoader" rql:"filter,sort"`
	QueryLoader  string `json:"queryLoader" rql:"filter,sort"`
	Loop         int    `json:"loop"`
	LoopInterval int    `json:"loopInterval" rql:"filter,sort"`

	Yy      string `gorm:"type:TEXT;" rql:"filter,sort"`
	Zz      string `gorm:"type:TEXT;" rql:"filter,sort"`
	Queries int    `json:"queries" rql:"filter,sort"`

	QueryFileName string
	DataFileName  string

	QueryStr string
	DataStr  string

	lock    sync.Mutex
	randgen *impl.RandgenLoader
}

//persistent the result and set the id
func (test *Test) Persistent() error {
	return db.GetDB().Create(test).Error
}

func (test *Test) Update() error {
	return db.GetDB().Save(test).Error
}

func GetTestById(id int64) (*Test, error) {
	result := &Test{}
	return result, db.GetDB().Where(" id=?", id).First(result).Error
}

func ListTest() ([]*Test, error) {
	var list []*Test
	return list, db.GetDB().Find(&list).Error
}

func (test *Test) GetName() string {
	return test.TestName
}

func (test *Test) GetDataLoaders() interfaces.DataLoader {
	switch test.DataLoader {
	case "file":
		return &impl.FileDataLoader{FileName: test.DataFileName}
	case "string":
		return &impl.StringLoader{SQLStr: test.DataStr}
	default:
		return test.getRandGen()
	}
}

func (test *Test) GetQueryLoaders() interfaces.QueryLoader {
	switch test.DataLoader {
	case "file":
		return &impl.FileDataLoader{FileName: test.QueryFileName}
	case "string":
		return &impl.StringLoader{SQLStr: test.QueryStr}
	default:
		return test.getRandGen()
	}
}

func (test *Test) GetComparor() interfaces.SqlResultComparer {
	return &sqldiff.StandardComparer{}
}

func (test *Test) GetLoop() int {
	return test.Loop
}

func (test *Test) GetLoopInterval() int {
	return test.LoopInterval
}

func (test *Test) getRandGen() *impl.RandgenLoader {
	test.lock.Lock()
	defer test.lock.Unlock()
	if test.randgen == nil {
		test.randgen = &impl.RandgenLoader{Yy: test.Yy, Zz: test.Zz, Queries: test.Queries}
	}
	return test.randgen
}
