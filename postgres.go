package migration

import (
	"fmt"
	"strings"
)

func generateAddColumnSQL(tableName string, field FieldInfo, aliasMap map[string]string) string {
	fieldType := field.Type
	if underlyingType, isAlias := aliasMap[fieldType]; isAlias {
		fieldType = underlyingType
	}
	sqlType := getSQLType(fieldType)
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, field.Name, sqlType)
}

func generateRemoveColumnSQL(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", tableName, columnName)
}

func generateAlterColumnTypeSQL(tableName, columnName, newType string) string {
	return fmt.Sprintf("-- ALTER TABLE %s ALTER COLUMN %s TYPE %s;", tableName, columnName, newType)
}

func generateCreateTableSQL(st StructInfo, aliasMap map[string]string) string {
	var columns []string
	for _, field := range st.Fields {
		columnName := field.Name
		columnType := field.Type

		if underlyingType, isAlias := aliasMap[columnType]; isAlias {
			columnType = underlyingType
		}

		sqlType := getSQLType(columnType)
		columns = append(columns, fmt.Sprintf("    %s %s", columnName, sqlType))
	}

	sql := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", st.Name, strings.Join(columns, ",\n"))
	return sql
}

func getSQLType(goType string) string {
	switch goType {
	case "string":
		return "TEXT"
	case "int", "int32":
		return "INTEGER"
	case "int64":
		return "BIGINT"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMP"
	case "[]byte":
		return "BYTEA"
	default:
		return "UNKNOWN"
	}
}
