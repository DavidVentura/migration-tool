package main

import (
	"log"

	"fmt"
	"github.com/davidventura/migration-tool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

func main() {
	dbConnectionString := "postgres://postgres:password@localhost:5555/postgres?sslmode=disable"

	db, err := sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	structs, aliases, err := migration.ScanCodebase("./")
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		os.Exit(2)
	}
	aliasMap := make(map[string]string)
	for _, alias := range aliases {
		aliasMap[alias.Name] = alias.UnderlyingType
	}

	tables := getDatabaseTables(db)

	mismatches := migration.CompareStructsAndTables(structs, aliasMap, tables)

	ignoredMismatches := 0
	for _, mismatch := range mismatches {
		if s, ok := mismatch.(migration.FieldMissingInStruct); ok {
			if s.Column.Name == "id" {
				ignoredMismatches += 1
				continue
			}
			if strings.Contains(s.Column.Name, "deprecated") {
				ignoredMismatches += 1
				continue
			}
		}

		fmt.Printf("%s\n", mismatch.Description())
		fmt.Printf("Suggested SQL:\n%s\n\n", mismatch.SQL())
	}

	if len(mismatches) == ignoredMismatches {
		fmt.Println("Schema matches codebase")
	} else {
		fmt.Println("Mismatches detected, exit 1")
		os.Exit(1)
	}
}
