package types

import (
	"github.com/monitoring-system/dbtest/db"
)

type LoopResult struct {
	ID          int64 `json:"id",gorm:"primary_key"`
	TestID      int64
	Loop        int
	Status      string
	Start       int64
	End         int64
	DML         string `gorm:"text"`
	Query       string `gorm:"text"`
	FailedDML   string `gorm:"text"`
	FailedQuery string `gorm:"text"`
}

func init() {
	db.GetDB().AutoMigrate(&LoopResult{})
}

//persistent the result and set the id
func (result *LoopResult) Persistent() error {
	return db.GetDB().Create(result).Error
}

func ListLoopResult(id int64) ([]*TestResult, error) {
	var list []*TestResult
	return list, db.GetDB().Where("test_id", id).Find(&list).Error
}
