package models

import (
	"binlog-db-sync/lib/parser"
	"encoding/json"
	"reflect"
	"strings"
)

type BaseModel struct {
	Properties interface{}
	columns    map[string]string
	tableName  string
	schemaName string
}

func (bm *BaseModel) TableName() string {
	return bm.tableName
}

func (bm *BaseModel) SchemaName() string {
	return bm.schemaName
}

func (bm *BaseModel) BuildJSON() string {
	jsonObject, _ := json.Marshal(bm.Properties)
	return replaceByMap(string(jsonObject), bm.columns)
}

type DBTable struct {
	SchemaName string
	TableName  string
}

func Create(properties interface{}, dbTable DBTable) (bm *BaseModel) {
	columns, _ := getColumns(properties)
	return &BaseModel{Properties: properties, columns: columns, schemaName: dbTable.SchemaName, tableName: dbTable.TableName}
}

func getColumns(properties interface{}) (columns map[string]string, err error) {
	var columnName string
	var ok bool

	columns = make(map[string]string)

	v := reflect.ValueOf(properties)
	s := reflect.Indirect(v)
	t := s.Type()
	num := t.NumField()
	for k := 0; k < num; k++ {
		parsedTag := parser.ParseTagSetting(t.Field(k).Tag)
		if columnName, ok = parsedTag["COLUMN"]; !ok || columnName == "COLUMN" {
			continue
		}
		columns[string(t.Field(k).Name)] = string(columnName)
	}

	return
}

func replaceByMap(s string, m map[string]string) string {
	for key, value := range m {
		s = strings.Replace(s, key, value, 1)
	}
	return s
}
