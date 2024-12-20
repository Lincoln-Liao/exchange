package usecase

import (
	"context"

	"exchange/internal/domain/transaction"
)

type TransactionUseCase struct {
	transactionService transaction.TransactionServiceInterface
}

func NewTransactionUseCase(tService transaction.TransactionServiceInterface) *TransactionUseCase {
	return &TransactionUseCase{
		transactionService: tService,
	}
}

func (uc *TransactionUseCase) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	txs, err := uc.transactionService.GetTransactionHistory(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func (uc *TransactionUseCase) GetTransactionByID(ctx context.Context, txID string) (transaction.Transaction, error) {
	tx, err := uc.transactionService.GetTransactionByID(ctx, txID)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return tx, nil
}
