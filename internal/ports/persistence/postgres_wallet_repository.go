package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"exchange/internal/domain/wallet"
)

type PostgresWalletRepository struct {
	db *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) *PostgresWalletRepository {
	return &PostgresWalletRepository{
		db: db,
	}
}

func (r *PostgresWalletRepository) CreateWallet(ctx context.Context, w wallet.Wallet) error {
	query := `
        INSERT INTO wallets (user_id, balance, currency, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(ctx, query,
		w.UserID, w.Balance, w.Currency, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *PostgresWalletRepository) GetWalletByUserID(ctx context.Context, userID string) (wallet.Wallet, error) {
	query := `
        SELECT user_id, balance, currency, created_at, updated_at
        FROM wallets
        WHERE user_id = $1
    `
	var w wallet.Wallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&w.UserID, &w.Balance, &w.Currency, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet.Wallet{}, wallet.ErrWalletNotFound
		}
		return wallet.Wallet{}, err
	}
	return w, nil
}

func (r *PostgresWalletRepository) UpdateWallet(ctx context.Context, w wallet.Wallet) error {
	query := `
        UPDATE wallets
        SET balance = $2, updated_at = $3
        WHERE user_id = $1
    `
	res, err := r.db.ExecContext(ctx, query, w.UserID, w.Balance, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return wallet.ErrWalletNotFound
	}

	return nil
}
