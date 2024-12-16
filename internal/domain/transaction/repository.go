package transaction

import (
	"context"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx Transaction) error

	GetTransactionByID(ctx context.Context, id string) (Transaction, error)

	ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]Transaction, error)
}
