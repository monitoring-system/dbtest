package sqldiff

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/filter"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/pingcap/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
)

type StandardComparer struct {
	Strict bool
}

func (c *StandardComparer) CompareQuery(db1, db2 *sql.DB, query string) (string, error, error) {
	expectedResult, err1 := GetQueryResult(db1, query)
	actualResult, err2 := GetQueryResult(db2, query)
	if err1 != nil || err2 != nil {
		return "", err1, err2
	} else {
		// now compare the results
		equals := ""
		if c.Strict {
			equals = c.strictCompare(expectedResult, actualResult)
		} else {
			equals = c.nonOrderCompare(expectedResult, actualResult)
		}
		return equals, nil, nil
	}
}

//return true if two result is equals with order
func (c *StandardComparer) strictCompare(expectedResult *SqlResult, actualResult *SqlResult) string {
	queryData1 := expectedResult.data
	queryData2 := actualResult.data
	if len(queryData1) != len(queryData2) {
		return fmt.Sprintf("expectedResult length not equals mysql=%d, tidb=%d", queryData1, len(queryData2))
	}

	for rowIndex, row := range queryData1 {
		if !c.compareRow(expectedResult.columnTypes, rowIndex, row, queryData2[rowIndex]) {
			return GetColorDiff(expectedResult.String(), actualResult.String())
		}
	}
	return ""
}

// compare two query result without ordered
func (c *StandardComparer) nonOrderCompare(result *SqlResult, result2 *SqlResult) string {
	queryData1 := result.data
	queryData2 := result2.data
	if len(queryData1) != len(queryData2) {
		return fmt.Sprintf("result length not equals mysql=%d,tidb=%d", len(queryData1), len(queryData2))
	}
	var checkedRowArray = make([]bool, len(queryData1))
	for rowIndex, row := range queryData1 {
		hasOneEquals := false
		for checkIndex, checked := range checkedRowArray {
			if !checked {
				equals := c.compareRow(result.columnTypes, rowIndex, row, queryData2[checkIndex])
				if equals {
					checkedRowArray[checkIndex] = true
					hasOneEquals = true
					break
				}
			}
		}
		if !hasOneEquals {
			return GetColorDiff(result.String(), result2.String())
		}
	}
	return ""
}

// compare two result row
func (c *StandardComparer) compareRow(columnTypes []*sql.ColumnType, rowIndex int, row [][]byte, row2 [][]byte) bool {
	//var line string
	for colIndex, col := range row {
		if len(row) != len(row2) {
			log.Info("result column length not equals", zap.Int("db1", len(row)), zap.Int("db2", len(row2)))
			return false
		}
		return c.compareCell(col, row2[colIndex], columnTypes[colIndex])
	}
	return true
}

func (c *StandardComparer) compareCell(cell1 []byte, cell2 []byte, columnType *sql.ColumnType) bool {
	cv1 := string(cell1)
	cv2 := string(cell2)

	//driver not support column type
	if cv1 != cv2 {
		//maybe it's json
		if (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) || (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) {
			if jsonEquals(cv1, cv2) {
				log.Info("result json equals", zap.String("cv1", cv1), zap.String("cv2", cv2))
				return true
			}
		}

		//now check if there is a custom cell compare
		return filter.FilterCompareDiff(cell1, cell2, columnType)
	}
	return true
}

func jsonEquals(s1, s2 string) bool {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(o1, o2)
}

func GetColorDiff(expect, actual string) string {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	patch := diffmatchpatch.New()
	diff := patch.DiffMain(expect, actual, false)
	var newExpectedContent, newActualResult bytes.Buffer
	for _, d := range diff {
		switch d.Type {
		case diffmatchpatch.DiffEqual:
			newExpectedContent.WriteString(d.Text)
			newActualResult.WriteString(d.Text)
		case diffmatchpatch.DiffDelete:
			newExpectedContent.WriteString(red(d.Text))
		case diffmatchpatch.DiffInsert:
			newActualResult.WriteString(green(d.Text))
		}
	}
	return fmt.Sprintf("Expected Result:\n%s\nActual Result:\n%s\n", newExpectedContent.String(), newActualResult.String())
}
