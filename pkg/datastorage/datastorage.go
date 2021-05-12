package datastorage

import "context"

// Inserter inserts records from a CSV file to a datasource.
type Inserter interface {
	Insert(ctx context.Context, records [][]string, options InsertOptions) error
}

// InsertOptions gathers the options of an insert process.
type InsertOptions struct {
	// Formatter formats each value of a record.
	// Optional.
	Formatter            func(value string) string
	// IgnoreDuplicateError ignores duplicate error from the datasource if set
	// to true.
	IgnoreDuplicateError bool
}
