package datastorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLInserter implements Inserter for a MySQL datasource.
type MySQLInserter struct {
	db *sql.DB

	table   string
	columns []interface{}
}

// NewMySQLInserter creates an instance of MySQLInserter.
func NewMySQLInserter(dsn, table string, columns []string, timeout time.Duration) (*MySQLInserter, error) {
	if table == "" {
		return nil, errors.New("empty table name")
	}

	if len(columns) == 0 {
		return nil, errors.New("no columns")
	}

	cols := make([]interface{}, len(columns))
	for i, col := range columns {
		cols[i] = col
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open connection to MySQL: %v", err)
	}

	db.SetConnMaxLifetime(timeout)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return &MySQLInserter{
		db:      db,
		table:   table,
		columns: cols,
	}, nil
}

// Insert inserts records to a MySQL table.
func (inserter *MySQLInserter) Insert(ctx context.Context, records [][]string, options InsertOptions) error {
	dialect := goqu.Dialect("mysql")
	builder := dialect.Insert(inserter.table).Cols(inserter.columns...)

	// Dummy formatter.
	format := func(value string) string {
		return value
	}

	if options.Formatter != nil {
		format = options.Formatter
	}

	for _, record := range records {
		values := make([]interface{}, len(record))

		for i, value := range record {
			values[i] = format(value)
		}

		builder = builder.Vals(values)
	}

	query, _, err := builder.ToSQL()
	if err != nil {
		return fmt.Errorf("build insert query: %v", err)
	}

	_, err = inserter.db.ExecContext(ctx, query)

	if err == nil {
		return nil
	}

	if options.IgnoreDuplicateError && strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
		return nil
	}

	return fmt.Errorf("execute insert query: %v", err)
}

// Interface guards.
var (
	_ Inserter = (*MySQLInserter)(nil)
)
