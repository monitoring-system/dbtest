package parser

import (
	"fmt"
	_ "fmt"
	"github.com/xwb1989/sqlparser"
	"reflect"
	"strings"
)

type Result struct {
	IsDDL     bool
	IgnoreSql bool
	TableName []string
	Rewrite   bool
	NewSql    string
}

func Parse(sql string) (*Result, error) {
	ast, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	switch ast.(type) {
	case *sqlparser.DDL:
		needRewrite, newSql := rewriteSql(ast.(*sqlparser.DDL))
		return &Result{
			IsDDL:     true,
			IgnoreSql: false,
			TableName: []string{getTbleNameByDDL(ast.(*sqlparser.DDL))},
			Rewrite:   needRewrite,
			NewSql:    newSql,
		}, nil
	case *sqlparser.Update, *sqlparser.Select:
		var tables []string
		return buildResult(false, getTableNames(reflect.Indirect(reflect.ValueOf(ast)), tables, 0, false)), nil
	case *sqlparser.Insert:
		in := ast.(*sqlparser.Insert)
		if in.Ignore == "" {
			return buildResult(false, []string{ast.(*sqlparser.Insert).Table.Name.String()}), nil
		} else {
			return buildResultWithIgnoreField(false, true, []string{ast.(*sqlparser.Insert).Table.Name.String()}), nil
		}
	}

	return buildResult(false, nil), nil
}

func getTbleNameByDDL(ddl *sqlparser.DDL) string {
	if ddl.NewName.IsEmpty() && ddl.Table.IsEmpty() {
		fmt.Println("parser table failed")
		return ""
	}

	if ddl.NewName.IsEmpty() {
		return ddl.Table.Name.String()
	}

	return ddl.NewName.Name.String()
}

func buildResult(isddl bool, tableNames []string) *Result {
	return buildResultWithIgnoreField(isddl, false, tableNames)
}

func buildResultWithIgnoreField(isddl bool, ignoresql bool, tableNames []string) *Result {
	return &Result{
		IsDDL:     isddl,
		IgnoreSql: ignoresql,
		TableName: tableNames,
	}
}

func getTableNames(v reflect.Value, tables []string, level int, isTable bool) []string {
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == "TableIdent" {
			// if this is a TableIdent struct, extract the table name
			tableName := v.FieldByName("v").String()
			if tableName != "" && isTable {
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

func rewriteSql(ddl *sqlparser.DDL) (needRewrite bool, sql string) {
	if ddl.TableSpec == nil || ddl.TableSpec.Columns == nil {
		return false, ""
	}

	for i, cd := range ddl.TableSpec.Columns {
		if strings.ToLower(cd.Type.Type) == "decimal" && cd.Type.Length == nil {
			ddl.TableSpec.Columns[i] = &sqlparser.ColumnDefinition{
				Name: cd.Name,
				Type: sqlparser.ColumnType{
					Type:          cd.Type.Type,
					NotNull:       cd.Type.NotNull,
					Autoincrement: cd.Type.Autoincrement,
					Default:       cd.Type.Default,
					OnUpdate:      cd.Type.OnUpdate,
					Comment:       cd.Type.Comment,
					Unsigned:      cd.Type.Unsigned,
					Zerofill:      cd.Type.Zerofill,
					Scale:         cd.Type.Scale,
					Charset:       cd.Type.Charset,
					Collate:       cd.Type.Collate,
					EnumValues:    cd.Type.EnumValues,
					KeyOpt:        cd.Type.KeyOpt,
					Length: &sqlparser.SQLVal{
						Type: sqlparser.IntVal,
						Val:  []byte("10"),
					},
				},
			}

			needRewrite = true
		}
	}

	if needRewrite {
		sql = sqlparser.String(ddl)
	}

	return
}
