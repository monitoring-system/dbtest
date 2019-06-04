package sqldiff

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/pingcap/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

// sql result present the result receive from database
type SqlResult struct {
	data        [][][]byte
	header      []string
	columnTypes []*sql.ColumnType
}

// readable query result like mysql shell client
func (result *SqlResult) String() string {
	if result.data == nil || result.header == nil {
		return "no result"
	}

	// Calculate the max column length
	var colLength []int
	for _, c := range result.header {
		colLength = append(colLength, len(c))
	}
	for _, row := range result.data {
		for n, col := range row {
			if l := len(col); colLength[n] < l {
				colLength[n] = l
			}
		}
	}
	// The total length
	var total = len(result.header) - 1
	for index := range colLength {
		colLength[index] += 2 // Value will wrap with space
		total += colLength[index]
	}

	var lines []string
	var push = func(line string) {
		lines = append(lines, line)
	}

	// Write table header
	var header string
	for index, col := range result.header {
		length := colLength[index]
		padding := length - 1 - len(col)
		if index == 0 {
			header += "|"
		}
		header += " " + col + strings.Repeat(" ", padding) + "|"
	}
	splitLine := "+" + strings.Repeat("-", total) + "+"
	push(splitLine)
	push(header)
	push(splitLine)

	// Write rows data
	for _, row := range result.data {
		var line string
		for index, col := range row {
			length := colLength[index]
			padding := length - 1 - len(col)
			if index == 0 {
				line += "|"
			}
			line += " " + string(col) + strings.Repeat(" ", padding) + "|"
		}
		push(line)
	}
	push(splitLine)
	return strings.Join(lines, "\n")
}

//return true if two result is equals with order
func (result *SqlResult) strictCompare(tidbResult *SqlResult) bool {
	queryData1 := result.data
	queryData2 := tidbResult.data
	if len(queryData1) != len(queryData2) {
		log.Info("result length not equals", zap.Int("db1", len(queryData1)), zap.Int("db2", len(queryData2)))
		return false
	}

	for rowIndex, row := range queryData1 {
		if !compareRow(rowIndex, row, queryData2[rowIndex]) {
			printColorDiff(result.String(), tidbResult.String())
			return false
		}
	}
	return true
}

// compare two query result without ordered
func (result *SqlResult) nonOrderCompare(result2 *SqlResult) bool {
	queryData1 := result.data
	queryData2 := result2.data
	if len(queryData1) != len(queryData2) {
		log.Info("result length not equals", zap.Int("db1", len(queryData1)), zap.Int("db2", len(queryData2)))
		return false
	}
	var checkedRowArray = make([]bool, len(queryData1))
	for rowIndex, row := range queryData1 {
		hasOneEquals := false
		for checkIndex, checked := range checkedRowArray {
			if !checked {
				equals := compareRow(rowIndex, row, queryData2[checkIndex])
				if equals {
					checkedRowArray[checkIndex] = true
					hasOneEquals = true
					break
				}
			}
		}
		if !hasOneEquals {
			printColorDiff(result.String(), result2.String())
			return false
		}
	}
	return true
}

// compare two result row
func compareRow(rowIndex int, row [][]byte, row2 [][]byte) bool {
	//var line string
	for colIndex, col := range row {
		if len(row) != len(row2) {
			log.Info("result column length not equals", zap.Int("db1", len(row)), zap.Int("db2", len(row2)))
			return false
		}

		cv1 := string(col)
		cv2 := string(row2[colIndex])
		//driver not support column type
		if cv1 != cv2 {
			//maybe it's json
			if (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) || (strings.HasPrefix(cv1, "{") && strings.HasPrefix(cv1, "{")) {
				if !jsonEquals(cv1, cv2) {
					log.Info("result json value not equals", zap.Int("row", rowIndex+1), zap.Int("col", colIndex+1), zap.String("cv1", cv1), zap.String("cv2", cv2))
					return false
				}
			} else {
				return false
			}
		}
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

func printColorDiff(expect, actual string) {
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
	fmt.Printf("Expected Result:\n%s\nActual Result:\n%s\n", newExpectedContent.String(), newActualResult.String())
}
