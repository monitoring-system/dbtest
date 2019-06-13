package types

import (
	"github.com/monitoring-system/dbtest/db"
)

const (
	TestStatusPending string = "pending"
	TestStatusRunning string = "running"
	TestStatusDone    string = "done"
	TestStatusFail    string = "fail"
)

func init() {
	db.GetDB().AutoMigrate(&TestResult{})
}

type TestResult struct {
	ID     int64 `json:"id",gorm:"primary_key"`
	TestID int64
	Name   string
	Status string
	Loop   int
}

//persistent the result and set the id
func (result *TestResult) Persistent() error {
	return db.GetDB().Create(result).Error
}

func (result *TestResult) Update() error {
	return db.GetDB().Save(result).Error
}

func GetTestResultByTestId(id int64) (*TestResult, error) {
	result := &TestResult{}
	return result, db.GetDB().Where("test_id", id).First(result).Error
}

func ListResult() ([]*TestResult, error) {
	var list []*TestResult
	return list, db.GetDB().Find(&list).Error
}
