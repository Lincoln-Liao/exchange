package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresTransactionManager struct {
	db *sql.DB
}

func NewPostgresTransactionManager(db *sql.DB) *PostgresTransactionManager {
	return &PostgresTransactionManager{db: db}
}

func (tm *PostgresTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctxTx := context.WithValue(ctx, transactionContextKey{}, tx)

	err = fn(ctxTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction after fn error: %v, original err: %w", rbErr, err)
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type transactionContextKey struct{}

func GetTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	v := ctx.Value(transactionContextKey{})
	if v == nil {
		return nil, false
	}
	tx, ok := v.(*sql.Tx)
	return tx, ok
}
