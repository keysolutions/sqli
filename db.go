package sqli

import (
	"context"
	"database/sql"
)

// DB is a wrapper around the database/sql database handle to provide the
// query scanner interface functions.
type DB struct {
	*sql.DB
}

// Row is an interface representing a row in the database that provides a Scan method.
// This interface is primiarly intended to represent a *sql.Row or *sql.Rows from the
// database/sql package.
type Row interface {
	Scan(dest ...interface{}) error
}

// Scanner is an interface to be used in scanning a database row handle.
type Scanner interface {
	Scan(row Row) error
}

// ScannerFn satifies the Scanner interface for a function type. The function will
// be called for each row.
type ScannerFn func(Row) error

// Scan calls the scanner function and returns its result.
func (s ScannerFn) Scan(row Row) error { return s(row) }

// Open opens a database specified by the database driver name and data source name.
// The arguments are passed to the database/sql Open function to retrive the
// database handle it provides.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Query provides a wrapper around the database/sql Query function which provides
// boilerplate around looping through the returned rows. The Scanner Scan method from
// the provided scanner is called on each interation of the row cursor loop.
//
// Query calls Scan on each row, which can be used like:
//
//  func (p *peopleScanner) Scan(row sqli.Row) error {
//      var person Person
//      if err := row.Scan(&person.Name); err != nil {
//           return err
//      }
//      p.people = append(p.people, &person)
//      return nil
//  }
func (db *DB) Query(scanner Scanner, query string, args ...interface{}) error {
	return db.QueryContext(context.Background(), scanner, query, args...)
}

// QueryIn operates like Query but expands arguments for IN clauses.
// IN arguments are to be passed as slices (e.g. []int{1, 2, 3}).
func (db *DB) QueryIn(scanner Scanner, query string, args ...interface{}) error {
	return db.QueryInContext(context.Background(), scanner, query, args...)
}

// QueryRow operates like Query but on a single row. The Scanner Scan function will
// be called only once.
func (db *DB) QueryRow(scanner Scanner, query string, args ...interface{}) error {
	return db.QueryRowContext(context.Background(), scanner, query, args...)
}

// QueryContext performs a query utilizing the supplied context.
func (db *DB) QueryContext(ctx context.Context, scanner Scanner, query string, args ...interface{}) error {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := scanner.Scan(rows)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return rows.Close()
}

// QueryInContext performs an in query utilizing the supplied context.
func (db *DB) QueryInContext(ctx context.Context, scanner Scanner, query string, args ...interface{}) error {
	query, args = In(query, args...)
	return db.QueryContext(ctx, scanner, query, args...)
}

// QueryRowContext performs a query on a single row utlizing the supplied context.
func (db *DB) QueryRowContext(ctx context.Context, scanner Scanner, query string, args ...interface{}) error {
	return scanner.Scan(db.DB.QueryRowContext(ctx, query, args...))
}
