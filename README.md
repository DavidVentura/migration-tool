# Database-Struct Comparison Tool

A Go tool for comparing Go structs with database tables, identifying mismatches, and generating SQL to resolve discrepancies.

## Features

- Scans Go codebase for structs annotated with `go:generate` containing the word "table"
- Compares struct fields with corresponding database table columns
- Identifies missing tables, missing fields, and type mismatches
- Generates SQL statements for creating tables, adding/removing columns, and altering column types
- Handles type aliases in Go structs

## Usage

1. Set up your database connection string in the `main` function.
2. Run the tool in your project root:
3. Review the output for mismatches and suggested SQL statements.

## Mismatch Types

- `TableMissing`: Struct without a corresponding database table
- `FieldMissingInTable`: Struct field missing from the database table
- `FieldMissingInStruct`: Database column missing from the struct
- `TypeMismatch`: Type discrepancy between struct field and database column

## Limitations

- Only supports postgres
- Doesn't handle complex scenarios like constraints or indexes


## License

[MIT License](https://opensource.org/licenses/MIT)
