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
	DML         string `gorm:"type:TEXT;"`
	Query       string `gorm:"type:TEXT;"`
	FailedDML   string `gorm:"type:TEXT;"`
	FailedQuery string `gorm:"type:TEXT;"`
}

//persistent the result and set the id
func (result *LoopResult) Persistent() error {
	return db.GetDB().Create(result).Error
}

func ListLoopResult(id int64) ([]*LoopResult, error) {
	var list []*LoopResult
	return list, db.GetDB().Where("test_id", id).Find(&list).Error
}
