package persistence

import (
	"context"
	"database/sql"
	"errors"

	"exchange/internal/domain/transaction"
)

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		db: db,
	}
}

func (r *PostgresTransactionRepository) CreateTransaction(ctx context.Context, tx transaction.Transaction) error {
	query := `
        INSERT INTO transactions (id, from_user_id, to_user_id, amount, currency, type, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query,
		tx.ID, tx.FromUserID, tx.ToUserID, tx.Amount, tx.Currency, string(tx.Type), tx.CreatedAt,
	)
	return err
}

func (r *PostgresTransactionRepository) GetTransactionByID(ctx context.Context, id string) (transaction.Transaction, error) {
	query := `
        SELECT id, from_user_id, to_user_id, amount, currency, type, created_at
        FROM transactions
        WHERE id = $1
    `
	var tx transaction.Transaction
	var tType string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Currency, &tType, &tx.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transaction.Transaction{}, transaction.ErrTransactionNotFound
		}
		return transaction.Transaction{}, err
	}
	tx.Type = transaction.TransactionType(tType)
	return tx, nil
}

func (r *PostgresTransactionRepository) ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	query := `
        SELECT id, from_user_id, to_user_id, amount, currency, type, created_at
        FROM transactions
        WHERE from_user_id = $1 OR to_user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []transaction.Transaction
	for rows.Next() {
		var tx transaction.Transaction
		var tType string
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.Currency, &tType, &tx.CreatedAt); err != nil {
			return nil, err
		}
		tx.Type = transaction.TransactionType(tType)
		results = append(results, tx)
	}

	return results, rows.Err()
}
