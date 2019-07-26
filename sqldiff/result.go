package sqldiff

import (
	"database/sql"
	"strings"
)

// sql result present the result receive from database
type SqlResult struct {
	data        [][][]byte
	header      []string
	columnTypes []*sql.ColumnType
}

func (result *SqlResult) HasResult() bool {
	if result.data == nil || result.header == nil {
		return false
	}
	return true
}

// draw readable query result like mysql shell client
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
