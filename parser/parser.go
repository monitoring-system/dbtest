package parser

import (
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
	"reflect"
)

type Result struct {
	IsDDL bool
	TableName []string
}

func Parse(sql string) (*Result, error) {
	ast, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	switch ast.(type) {
	case *sqlparser.DDL:
		return &Result{
			IsDDL:true,
		TableName: []string{ast.(*sqlparser.DDL).Table.Name.String()},
	}, nil
	case *sqlparser.Update,*sqlparser.Select:
		var tables []string
		tables = getTableNames(reflect.Indirect(reflect.ValueOf(ast)), tables, 0, false)
		return &Result{
			IsDDL: false,
			TableName: tables,
		}, nil
	case *sqlparser.Insert:
		return &Result{
			IsDDL: false,
			TableName: []string{ast.(*sqlparser.Insert).Table.Name.String()},
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("Unsupport statement, type=%T", ast))
}

func getTableNames(v reflect.Value, tables []string, level int, isTable bool) []string {
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == "TableIdent" {
			// if this is a TableIdent struct, extract the table name
			tableName := v.FieldByName("v").String()
			if tableName != "" && isTable{
				tables = append(tables, tableName)
			}
		} else {
			// otherwise enumerate all fields of the struct and process further
			for i := 0; i < v.NumField(); i++ {
				tables = getTableNames(reflect.Indirect(v.Field(i)), tables, level+1, isTable)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			// enumerate all elements of an array/slice and process further
			tables = getTableNames(reflect.Indirect(v.Index(i)), tables, level+1, isTable)
		}
	case reflect.Interface:
		if v.Type().Name() == "SimpleTableExpr" {
			isTable = true
		}
		// get the actual object that satisfies an interface and process further
		tables = getTableNames(reflect.Indirect(reflect.ValueOf(v.Interface())), tables, level+1, isTable)
	}

	return tables
}
