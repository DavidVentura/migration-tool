package migration

import (
	"reflect"
	"testing"
)

func TestCompareStructsAndTables(t *testing.T) {
	// Mock Go code as a string
	goCode := []byte(`
package model

//go:generate table
type User struct {
	ID        int    ` + "`json:\"id\"`" + `
	Name      string ` + "`json:\"name\"`" + `
	Email     string ` + "`json:\"email\"`" + `
	CreatedAt Time   ` + "`json:\"created_at\"`" + `
}

type Time string
`)

	// Parse the Go code to get StructInfo
	structs, aliases, err := parseFile("test.go", goCode)
	if err != nil {
		t.Fatalf("Failed to parse Go code: %v", err)
	}

	// Create a map of aliases
	aliasMap := make(map[string]string)
	for _, alias := range aliases {
		aliasMap[alias.Name] = alias.UnderlyingType
	}

	// Mock TableInfo
	tables := []TableInfo{
		{
			Name: "User",
			Columns: []ColumnInfo{
				{Name: "ID", Type: "INTEGER"},
				{Name: "Name", Type: "TEXT"},
				{Name: "Email", Type: "TEXT"},
				// CreatedAt is missing intentionally to create a mismatch
				{Name: "UpdatedAt", Type: "TIMESTAMP"}, // This column is not in the struct
			},
		},
	}

	mismatches := CompareStructsAndTables(structs, aliasMap, tables)

	// Expected mismatches
	expectedMismatches := []Mismatch{
		FieldMissingInTable{Table: "User", Field: FieldInfo{Name: "CreatedAt", Type: "Time"}},
		FieldMissingInStruct{Table: "User", Column: ColumnInfo{Name: "UpdatedAt", Type: "TIMESTAMP"}},
	}

	// Compare the results
	if len(mismatches) != len(expectedMismatches) {
		t.Errorf("Expected %d mismatches, but got %d", len(expectedMismatches), len(mismatches))
	}

	for i, mismatch := range mismatches {
		if reflect.TypeOf(mismatch) != reflect.TypeOf(expectedMismatches[i]) {
			t.Errorf("Mismatch type at index %d doesn't match. Expected %T, got %T", i, expectedMismatches[i], mismatch)
		}

		// You might want to add more specific checks here, depending on your Mismatch implementations
		if mismatch.Description() != expectedMismatches[i].Description() {
			t.Errorf("Mismatch description at index %d doesn't match. Expected %s, got %s", i, expectedMismatches[i].Description(), mismatch.Description())
		}
	}
}
