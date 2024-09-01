package migration

import (
	"strings"
)

func CompareStructsAndTables(structs []StructInfo, aliasMap map[string]string, tables []TableInfo) []Mismatch {
	var mismatches []Mismatch
	for _, structInfo := range structs {
		for i, field := range structInfo.Fields {
			if underlyingType, isAlias := aliasMap[field.Type]; isAlias {
				structInfo.Fields[i].Type = underlyingType
				structInfo.Fields[i].IsAlias = true
			}
		}
	}
	for _, st := range structs {
		found := false
		for _, tb := range tables {
			if !strings.EqualFold(st.Name, tb.Name) {
				continue
			}
			found = true
			mismatches = append(mismatches, compareFieldsAndColumns(st, tb, aliasMap)...)
			break
		}
		if !found {
			mismatches = append(mismatches, TableMissing{Name: st.Name, Fields: st.Fields})
		}
	}
	return mismatches
}

func compareFieldsAndColumns(st StructInfo, tb TableInfo, aliasMap map[string]string) []Mismatch {
	var mismatches []Mismatch

	fieldMap := make(map[string]FieldInfo)
	for _, f := range st.Fields {
		fieldMap[strings.ToLower(f.Name)] = f
	}

	columnMap := make(map[string]ColumnInfo)
	for _, c := range tb.Columns {
		columnMap[strings.ToLower(c.Name)] = c
	}

	for name, field := range fieldMap {
		if col, ok := columnMap[name]; ok {
			fieldType := field.Type
			if underlyingType, isAlias := aliasMap[fieldType]; isAlias {
				fieldType = underlyingType
			}
			if !compareTypes(field, col.Type) {
				mismatches = append(mismatches, TypeMismatch{
					Table:      tb.Name,
					FieldName:  field.Name,
					StructType: fieldType,
					TableType:  col.Type,
				})
			}
			delete(columnMap, name)
		} else {
			mismatches = append(mismatches, FieldMissingInTable{Table: tb.Name, Field: field})
		}
	}

	for _, col := range columnMap {
		mismatches = append(mismatches, FieldMissingInStruct{Table: tb.Name, Column: col})
	}

	return mismatches
}
func compareTypes(field FieldInfo, dbType string) bool {
	structType := field.Type
	dbType = strings.ToLower(dbType)

	switch dbType {
	case "text", "character varying", "char":
		return structType == "string"
	case "integer", "bigint":
		return structType == "int" || structType == "int32" || structType == "int64"
	case "real", "double precision":
		return structType == "float32" || structType == "float64"
	case "boolean":
		return structType == "bool"
	case "timestamp", "date":
		return structType == "time.time"
	case "numeric", "decimal":
		return structType == "float64" || structType == "decimal.Decimal"
	case "bytea":
		return structType == "[]byte"
	case "json":
		return structType == "map[string]interface{}" || structType == "interface{}"
	case "uuid":
		return structType == "string" // shouldn't this be a UUID package
	case "timestamp with time zone", "timestamp without time zone":
		return structType == "time.Time"
	case "inet":
		// should be net.IP no?
		return structType == "string"
	case "jsonb":
		return structType == "[]byte"
	case "user-defined":
		// any enum type in the DB should be an alias in code
		return field.IsAlias
	default:
		return false
	}
}
