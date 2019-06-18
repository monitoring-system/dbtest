package types

import (
	"github.com/a8m/rql"
	"github.com/monitoring-system/dbtest/db"
)

const (
	TestStatusPending string = "pending"
	TestStatusRunning string = "running"
	TestStatusDone    string = "done"

	TestStatusOK   string = "OK"
	TestStatusSkip string = "skip"
	TestStatusFail string = "fail"
)

type TestResult struct {
	ID              int64  `json:"id" gorm:"primary_key" rql:"filter,sort"`
	TestID          int64  `rql:"filter,sort"`
	Name            string `rql:"filter,sort"`
	Status          string `rql:"filter,sort"`
	FailedLoopCount int    `rql:"filter,sort"`
	Loop            int    `rql:"filter,sort"`
	Start           int64  `rql:"filter,sort"`
	End             int64  `rql:"filter,sort"`
}

//persistent the result and set the id
func (result *TestResult) Persistent() error {
	return db.GetDB().Create(result).Error
}

func (result *TestResult) Update() error {
	return db.GetDB().Save(result).Error
}

func ListUnFinishedTestResult() ([]*TestResult, error) {
	var list []*TestResult
	return list, db.GetDB().Where(" status != ?", TestStatusDone).Find(&list).Error
}

func ListResult(p *rql.Params) ([]*TestResult, error) {
	var list []*TestResult
	return list, db.GetDB().Where(p.FilterExp, p.FilterArgs).
		Offset(p.Offset).
		Limit(p.Limit).
		Order(p.Sort).Find(&list).Error
}
