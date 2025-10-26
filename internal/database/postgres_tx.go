package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// PostgresTx wraps a pgx.Tx to implement our database-agnostic Tx interface
type PostgresTx struct {
	tx pgx.Tx
}

// Commit commits the transaction
func (p *PostgresTx) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

// Rollback rolls back the transaction
func (p *PostgresTx) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}

// Unwrap returns the underlying pgx.Tx for PostgreSQL-specific operations
// This is used internally by the PostgreSQL database implementation
func (p *PostgresTx) Unwrap() pgx.Tx {
	return p.tx
}
