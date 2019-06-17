package types

import (
	"github.com/monitoring-system/dbtest/db"
)

type LoopResult struct {
	ID     int64  `json:"id",gorm:"primary_key"`
	TestID int64  `rql:"filter,sort"`
	Loop   int    `rql:"filter,sort"`
	Status string `rql:"filter,sort"`
	Start  int64  `rql:"filter,sort"`
	End    int64  `rql:"filter,sort"`
	//DML    string `gorm:"type:longtext;"`
	//Query  string `gorm:"type:longtext;"`
}

//persistent the result and set the id
func (result *LoopResult) Persistent() error {
	return db.GetDB().Create(result).Error
}

func ListLoopResult(id int64) ([]*LoopResult, error) {
	var list []*LoopResult
	return list, db.GetDB().Where("test_id", id).Find(&list).Error
}
