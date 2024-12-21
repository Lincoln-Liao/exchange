package transaction

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
)

type TransactionServiceInterface interface {
	LogTransaction(ctx context.Context, fromUserID, toUserID string, amount int64, currency string, tType TransactionType) (Transaction, error)
	GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (Transaction, error)
}

type TransactionService struct {
	repository TransactionRepository
}

func NewTransactionService(repo TransactionRepository) *TransactionService {
	return &TransactionService{
		repository: repo,
	}
}

func (s *TransactionService) LogTransaction(ctx context.Context, fromUserID, toUserID string, amount int64, currency string, tType TransactionType) (Transaction, error) {
	if amount <= 0 {
		return Transaction{}, ErrInvalidTransactionAmount
	}

	if tType != TransactionTypeDeposit && tType != TransactionTypeWithdraw && tType != TransactionTypeTransfer {
		return Transaction{}, ErrInvalidTransactionType
	}

	id, err := generateTransactionID()
	if err != nil {
		return Transaction{}, err
	}

	tx, err := NewTransaction(id, fromUserID, toUserID, amount, currency, tType)
	if err := s.repository.CreateTransaction(ctx, tx); err != nil {
		return Transaction{}, err
	}

	return tx, nil
}

func (s *TransactionService) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]Transaction, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	txs, err := s.repository.ListTransactionsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, ErrDatabaseFailure
	}

	return txs, nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id string) (Transaction, error) {
	if id == "" {
		return Transaction{}, ErrInvalidTransactionID
	}

	tx, err := s.repository.GetTransactionByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			return Transaction{}, ErrTransactionNotFound
		}
		return Transaction{}, err
	}
	return tx, nil
}

func generateTransactionID() (string, error) {
	id, err := uuid.NewV7()
	return id.String(), err
}
