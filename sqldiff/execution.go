package sqldiff

import "database/sql"

//execute the sql statement and return the raw data
func GetQueryResult(db *sql.DB, query string) (*SqlResult, error) {
	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	cols, err := result.Columns()
	if err != nil {
		return nil, err
	}
	types, err := result.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var allRows [][][]byte
	for result.Next() {
		var columns = make([][]byte, len(cols))
		var pointer = make([]interface{}, len(cols))
		for i := range columns {
			pointer[i] = &columns[i]
		}
		err := result.Scan(pointer...)
		if err != nil {
			return nil, err
		}
		allRows = append(allRows, columns)
	}
	queryResult := SqlResult{data: allRows, header: cols, columnTypes: types}
	return &queryResult, nil
}
