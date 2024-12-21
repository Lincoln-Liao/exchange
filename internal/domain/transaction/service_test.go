package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, tx Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]Transaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByID(ctx context.Context, id string) (Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Transaction), args.Error(1)
}

func TestTransactionService_LogTransaction(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	ctx := context.Background()

	t.Run("successful log transaction (deposit)", func(t *testing.T) {
		fromUserID := ""
		toUserID := "user1"
		amount := int64(1000)
		currency := "USD"
		tType := TransactionTypeDeposit

		expectedTx := Transaction{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Currency:   currency,
			Type:       tType,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx Transaction) bool {
			return tx.FromUserID == expectedTx.FromUserID &&
				tx.ToUserID == expectedTx.ToUserID &&
				tx.Amount == expectedTx.Amount &&
				tx.Currency == expectedTx.Currency &&
				tx.Type == expectedTx.Type
		})).Return(nil)

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.NoError(t, err)
		assert.Equal(t, fromUserID, tx.FromUserID)
		assert.Equal(t, toUserID, tx.ToUserID)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, currency, tx.Currency)
		assert.Equal(t, tType, tx.Type)
		assert.WithinDuration(t, time.Now(), tx.CreatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("successful log transaction (withdraw)", func(t *testing.T) {
		fromUserID := "user1"
		toUserID := ""
		amount := int64(500)
		currency := "USD"
		tType := TransactionTypeWithdraw

		expectedTx := Transaction{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Currency:   currency,
			Type:       tType,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx Transaction) bool {
			return tx.FromUserID == expectedTx.FromUserID &&
				tx.ToUserID == expectedTx.ToUserID &&
				tx.Amount == expectedTx.Amount &&
				tx.Currency == expectedTx.Currency &&
				tx.Type == expectedTx.Type
		})).Return(nil)

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.NoError(t, err)
		assert.Equal(t, fromUserID, tx.FromUserID)
		assert.Equal(t, toUserID, tx.ToUserID)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, currency, tx.Currency)
		assert.Equal(t, tType, tx.Type)
		assert.WithinDuration(t, time.Now(), tx.CreatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("successful log transaction (transfer)", func(t *testing.T) {
		fromUserID := "user1"
		toUserID := "user2"
		amount := int64(300)
		currency := "USD"
		tType := TransactionTypeTransfer

		expectedTx := Transaction{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Currency:   currency,
			Type:       tType,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx Transaction) bool {
			return tx.FromUserID == expectedTx.FromUserID &&
				tx.ToUserID == expectedTx.ToUserID &&
				tx.Amount == expectedTx.Amount &&
				tx.Currency == expectedTx.Currency &&
				tx.Type == expectedTx.Type
		})).Return(nil)

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.NoError(t, err)
		assert.Equal(t, fromUserID, tx.FromUserID)
		assert.Equal(t, toUserID, tx.ToUserID)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, currency, tx.Currency)
		assert.Equal(t, tType, tx.Type)
		assert.WithinDuration(t, time.Now(), tx.CreatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid amount", func(t *testing.T) {
		fromUserID := "user1"
		toUserID := "user2"
		amount := int64(-100)
		currency := "USD"
		tType := TransactionTypeTransfer

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.Error(t, err)
		assert.Equal(t, Transaction{}, tx)
		assert.EqualError(t, err, ErrInvalidTransactionAmount.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid transaction type", func(t *testing.T) {
		fromUserID := "user1"
		toUserID := "user2"
		amount := int64(100)
		currency := "USD"
		tType := TransactionType("INVALID_TYPE")

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.Error(t, err)
		assert.Equal(t, Transaction{}, tx)
		assert.EqualError(t, err, ErrInvalidTransactionType.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository create transaction failure", func(t *testing.T) {
		fromUserID := "user1"
		toUserID := "user2"
		amount := int64(100)
		currency := "USD"
		tType := TransactionTypeTransfer

		mockRepo.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tx Transaction) bool {
			return tx.FromUserID == fromUserID &&
				tx.ToUserID == toUserID &&
				tx.Amount == amount &&
				tx.Currency == currency &&
				tx.Type == tType
		})).Return(ErrDatabaseFailure)

		tx, err := service.LogTransaction(ctx, fromUserID, toUserID, amount, currency, tType)

		assert.Error(t, err)
		assert.Equal(t, Transaction{}, tx)
		assert.EqualError(t, err, ErrDatabaseFailure.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetTransactionHistory(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	ctx := context.Background()

	t.Run("successful get transaction history", func(t *testing.T) {
		userID := "user1"
		limit := 10
		offset := 0

		expectedTxs := []Transaction{
			{
				ID:         "tx1",
				FromUserID: "",
				ToUserID:   "user1",
				Amount:     1000,
				Currency:   "USD",
				Type:       TransactionTypeDeposit,
				CreatedAt:  time.Now(),
			},
			{
				ID:         "tx2",
				FromUserID: "user1",
				ToUserID:   "user2",
				Amount:     500,
				Currency:   "USD",
				Type:       TransactionTypeTransfer,
				CreatedAt:  time.Now(),
			},
		}

		mockRepo.On("ListTransactionsByUserID", ctx, userID, limit, offset).Return(expectedTxs, nil)

		txs, err := service.GetTransactionHistory(ctx, userID, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, expectedTxs, txs)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		userID := ""
		limit := 10
		offset := 0

		txs, err := service.GetTransactionHistory(ctx, userID, limit, offset)

		assert.Error(t, err)
		assert.Equal(t, []Transaction(nil), txs)
		assert.EqualError(t, err, ErrInvalidUserID.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetTransactionByID(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)

	ctx := context.Background()

	t.Run("successful get transaction by ID", func(t *testing.T) {
		id := "tx1"
		expectedTx := Transaction{
			ID:         id,
			FromUserID: "",
			ToUserID:   "user1",
			Amount:     1000,
			Currency:   "USD",
			Type:       TransactionTypeDeposit,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("GetTransactionByID", ctx, id).Return(expectedTx, nil)

		tx, err := service.GetTransactionByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedTx, tx)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid transaction ID", func(t *testing.T) {
		id := ""

		tx, err := service.GetTransactionByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, Transaction{}, tx)
		assert.EqualError(t, err, ErrInvalidTransactionID.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("transaction not found", func(t *testing.T) {
		id := "tx999"

		mockRepo.On("GetTransactionByID", ctx, id).Return(Transaction{}, ErrTransactionNotFound)

		tx, err := service.GetTransactionByID(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, Transaction{}, tx)
		assert.Equal(t, ErrTransactionNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}
