package main

import (
	"log"

	"github.com/davidventura/migration-tool"
	"github.com/jmoiron/sqlx"
)

func getDatabaseTables(db *sqlx.DB) []migration.TableInfo {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error querying tables: %v", err)
	}
	defer rows.Close()

	var tables []migration.TableInfo
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("Error scanning table name: %v", err)
			continue
		}

		columns := getTableColumns(db, tableName)
		tables = append(tables, migration.TableInfo{Name: tableName, Columns: columns})
	}

	return tables
}

func getTableColumns(db *sqlx.DB, tableName string) []migration.ColumnInfo {
	query := `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = $1
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		log.Fatalf("Error querying columns for table %s: %v", tableName, err)
	}
	defer rows.Close()

	var columns []migration.ColumnInfo
	for rows.Next() {
		var colName, colType string
		if err := rows.Scan(&colName, &colType); err != nil {
			log.Printf("Error scanning column info: %v", err)
			continue
		}
		columns = append(columns, migration.ColumnInfo{Name: colName, Type: colType})
	}

	return columns
}
