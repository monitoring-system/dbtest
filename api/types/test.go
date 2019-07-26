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
	Comparer     string `json:"comparer"`
	SqlLoader    string `json:"sqlLoader" rql:"filter,sort"`
	Loop         int    `json:"loop"`
	LoopInterval int    `json:"loopInterval" rql:"filter,sort"`

	Yy      string `gorm:"type:TEXT;" rql:"filter,sort"`
	Zz      string `gorm:"type:TEXT;" rql:"filter,sort"`
	Queries int    `json:"queries" rql:"filter,sort"`

	DbRetain bool

	// unuseful now
	QueryFileName string
	DataFileName  string

	QueryStr string
	DataStr  string

	StrictCompare bool

	ContainContent string

	lock         sync.Mutex
	randgen      *impl.RandgenLoader
}

//persistent the result and set the id
func (this *Test) Persistent() error {
	return db.GetDB().Create(this).Error
}

func (this *Test) Update() error {
	return db.GetDB().Save(this).Error
}

func GetTestById(id int64) (*Test, error) {
	result := &Test{}
	return result, db.GetDB().Where(" id=?", id).First(result).Error
}

func ListTest() ([]*Test, error) {
	var list []*Test
	return list, db.GetDB().Find(&list).Error
}

func (this *Test) GetName() string {
	return this.TestName
}

func (this *Test) GetSqlLoader() interfaces.SqlLoader {
	switch this.SqlLoader {
	case impl.String:
		return &impl.StringLoader{SQLStr: this.DataStr}
	default:
		return &impl.RandgenLoader{Yy: this.Yy, Zz: this.Zz, Queries: this.Queries}
	}
}

func (this *Test) GetComparor() interfaces.SqlResultComparer {
	switch this.Comparer {
	case impl.Contain:
		return &impl.ContainsComparer{Content: this.ContainContent}
	default:
		// always there now
		return &sqldiff.StandardComparer{Strict: this.StrictCompare}
	}
}

func (this *Test) GetLoop() int {
	return this.Loop
}

func (this *Test) GetLoopInterval() int {
	return this.LoopInterval
}

func (this *Test) getRandGen() *impl.RandgenLoader {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.randgen == nil {
		this.randgen = &impl.RandgenLoader{Yy: this.Yy, Zz: this.Zz, Queries: this.Queries}
	}
	return this.randgen
}
