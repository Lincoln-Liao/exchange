package usecase

import (
	"context"
	"testing"

	"exchange/internal/domain/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransactionByID(ctx context.Context, txID string) (transaction.Transaction, error) {
	args := m.Called(ctx, txID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionService) LogTransaction(ctx context.Context, fromUserID, toUserID string, amount int64, currency string, tType transaction.TransactionType) (transaction.Transaction, error) {
	args := m.Called(ctx, fromUserID, toUserID, amount, currency, tType)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func TestTransactionUseCase_GetTransactionHistory(t *testing.T) {
	mockService := new(MockTransactionService)
	useCase := NewTransactionUseCase(mockService)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		userID := "user1"
		limit := 10
		offset := 0

		expectedTxs := []transaction.Transaction{
			{
				ID:         "tx1",
				FromUserID: "user1",
				ToUserID:   "user2",
				Amount:     1000,
				Currency:   "USD",
				Type:       transaction.TransactionTypeDeposit,
			},
			{
				ID:         "tx2",
				FromUserID: "user1",
				ToUserID:   "user3",
				Amount:     2000,
				Currency:   "USD",
				Type:       transaction.TransactionTypeTransfer,
			},
		}

		mockService.On("GetTransactionHistory", ctx, userID, limit, offset).Return(expectedTxs, nil)

		txs, err := useCase.GetTransactionHistory(ctx, userID, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, expectedTxs, txs)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		userID := ""
		limit := 10
		offset := 0

		mockService.On("GetTransactionHistory", ctx, userID, limit, offset).Return([]transaction.Transaction(nil), transaction.ErrInvalidUserID)

		txs, err := useCase.GetTransactionHistory(ctx, userID, limit, offset)

		assert.ErrorIs(t, err, transaction.ErrInvalidUserID)
		assert.Equal(t, txs, []transaction.Transaction(nil))
		mockService.AssertExpectations(t)
	})

	t.Run("repository failure", func(t *testing.T) {
		userID := "user1"
		limit := -1
		offset := 0

		mockService.On("GetTransactionHistory", ctx, userID, limit, offset).Return([]transaction.Transaction(nil), transaction.ErrDatabaseFailure)

		txs, err := useCase.GetTransactionHistory(ctx, userID, limit, offset)

		assert.ErrorIs(t, err, transaction.ErrDatabaseFailure)
		assert.Equal(t, txs, []transaction.Transaction(nil))
		mockService.AssertExpectations(t)
	})
}

func TestTransactionUseCase_GetTransactionByID(t *testing.T) {
	mockService := new(MockTransactionService)
	useCase := NewTransactionUseCase(mockService)

	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		txID := "tx123"
		expectedTx := transaction.Transaction{
			ID:         txID,
			FromUserID: "user1",
			ToUserID:   "user2",
			Amount:     1000,
			Currency:   "USD",
			Type:       transaction.TransactionTypeDeposit,
		}

		mockService.On("GetTransactionByID", ctx, txID).Return(expectedTx, nil)

		tx, err := useCase.GetTransactionByID(ctx, txID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTx, tx)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid transaction ID", func(t *testing.T) {
		txID := ""

		mockService.On("GetTransactionByID", ctx, txID).Return(transaction.Transaction{}, transaction.ErrInvalidTransactionID)

		tx, err := useCase.GetTransactionByID(ctx, txID)

		assert.ErrorIs(t, err, transaction.ErrInvalidTransactionID)
		assert.Equal(t, transaction.Transaction{}, tx)
		mockService.AssertExpectations(t)
	})

	t.Run("transaction not found", func(t *testing.T) {
		txID := "nonexistent"

		mockService.On("GetTransactionByID", ctx, txID).Return(transaction.Transaction{}, transaction.ErrTransactionNotFound)

		tx, err := useCase.GetTransactionByID(ctx, txID)

		assert.ErrorIs(t, err, transaction.ErrTransactionNotFound)
		assert.Equal(t, transaction.Transaction{}, tx)
		mockService.AssertExpectations(t)
	})

}
