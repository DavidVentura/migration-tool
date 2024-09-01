package migration

import (
	"fmt"
)

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type FieldInfo struct {
	Name    string
	Type    string
	IsAlias bool
	Tag     string
}

type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}

type ColumnInfo struct {
	Name string
	Type string
}

type TypeAlias struct {
	Name           string
	UnderlyingType string
}
type Mismatch interface {
	Description() string
	SQL() string
}

type TableMissing struct {
	Name   string
	Fields []FieldInfo
}

func (tm TableMissing) Description() string {
	return fmt.Sprintf("Table missing: %s", tm.Name)
}

func (tm TableMissing) SQL() string {
	return generateCreateTableSQL(StructInfo{Name: tm.Name, Fields: tm.Fields}, nil)
}

type FieldMissingInTable struct {
	Table string
	Field FieldInfo
}

func (fm FieldMissingInTable) Description() string {
	return fmt.Sprintf("Field '%s' is in struct but not in table '%s'", fm.Field.Name, fm.Table)
}

func (fm FieldMissingInTable) SQL() string {
	return generateAddColumnSQL(fm.Table, fm.Field, nil)
}

type FieldMissingInStruct struct {
	Table  string
	Column ColumnInfo
}

func (fm FieldMissingInStruct) Description() string {
	return fmt.Sprintf("Field '%s' is in table '%s' but not in struct", fm.Column.Name, fm.Table)
}

func (fm FieldMissingInStruct) SQL() string {
	return generateRemoveColumnSQL(fm.Table, fm.Column.Name)
}

type TypeMismatch struct {
	Table      string
	FieldName  string
	StructType string
	TableType  string
}

func (tm TypeMismatch) Description() string {
	return fmt.Sprintf("Type mismatch in table %s for field %s: struct has %s, table has %s",
		tm.Table, tm.FieldName, tm.StructType, tm.TableType)
}

func (tm TypeMismatch) SQL() string {
	return generateAlterColumnTypeSQL(tm.Table, tm.FieldName, getSQLType(tm.StructType))
}
