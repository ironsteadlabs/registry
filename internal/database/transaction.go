package database

import (
	"context"
)

// Tx represents a database transaction
// This abstracts the underlying database transaction implementation
type Tx interface {
	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error
}

// Row represents a single row result from a database query
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows represents multiple row results from a database query
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
	Err() error
}
