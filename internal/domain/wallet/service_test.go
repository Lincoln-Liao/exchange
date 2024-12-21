package wallet

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, w Wallet) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletByUserID(ctx context.Context, userID string) (Wallet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWallet(ctx context.Context, w Wallet) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func TestWalletService_CreateNewWallet(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	currency := "USD"

	t.Run("successful wallet creation", func(t *testing.T) {
		expectedWallet := Wallet{
			UserID:    userID,
			Balance:   0,
			Currency:  currency,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("CreateWallet", ctx, mock.MatchedBy(func(w Wallet) bool {
			return w.UserID == expectedWallet.UserID &&
				w.Balance == expectedWallet.Balance &&
				w.Currency == expectedWallet.Currency
		})).Return(nil)

		w, err := service.CreateNewWallet(ctx, userID, currency)

		assert.NoError(t, err)
		assert.Equal(t, expectedWallet.UserID, w.UserID)
		assert.Equal(t, expectedWallet.Balance, w.Balance)
		assert.Equal(t, expectedWallet.Currency, w.Currency)
		assert.WithinDuration(t, time.Now(), w.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), w.UpdatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_Deposit(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	initialBalance := int64(1000)
	depositAmount := int64(500)
	currency := "USD"

	existingWallet := Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful deposit", func(t *testing.T) {
		mockRepo.On("GetWalletByUserID", ctx, userID).Return(existingWallet, nil)

		updatedWallet := existingWallet
		updatedWallet.AddBalance(depositAmount)

		mockRepo.On("UpdateWallet", ctx, mock.MatchedBy(func(w Wallet) bool {
			return w.Balance == initialBalance+depositAmount &&
				w.UserID == userID &&
				w.Currency == currency
		})).Return(nil)

		err := service.Deposit(ctx, userID, depositAmount)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid amount (zero)", func(t *testing.T) {
		err := service.Deposit(ctx, userID, 0)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAmount, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid amount (negative)", func(t *testing.T) {
		err := service.Deposit(ctx, userID, -100)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAmount, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockRepo.On("GetWalletByUserID", ctx, "userempty").Return(Wallet{}, ErrWalletNotFound)

		err := service.Deposit(ctx, "userempty", depositAmount)

		assert.Error(t, err)
		assert.Equal(t, ErrWalletNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestWalletService_Withdraw 測試 Withdraw 方法
func TestWalletService_Withdraw(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	initialBalance := int64(1000)
	withdrawAmount := int64(500)
	currency := "USD"

	existingWallet := Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful withdraw", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetWalletByUserID", ctx, userID).Return(existingWallet, nil)

		updatedWallet := existingWallet
		updatedWallet.SubtractBalance(withdrawAmount)

		mockRepo.On("UpdateWallet", ctx, mock.MatchedBy(func(w Wallet) bool {
			return w.Balance == initialBalance-withdrawAmount &&
				w.UserID == userID &&
				w.Currency == currency
		})).Return(nil)

		err := service.Withdraw(ctx, userID, withdrawAmount)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid amount (zero)", func(t *testing.T) {
		err := service.Withdraw(ctx, userID, 0)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAmount, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid amount (negative)", func(t *testing.T) {
		err := service.Withdraw(ctx, userID, -100)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidAmount, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockRepo.On("GetWalletByUserID", ctx, "userempty").Return(Wallet{}, ErrWalletNotFound)

		err := service.Withdraw(ctx, "userempty", withdrawAmount)

		assert.Error(t, err)
		assert.Equal(t, ErrWalletNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		withdrawAmount := int64(1500)

		mockRepo.On("GetWalletByUserID", ctx, userID).Return(existingWallet, nil)

		err := service.Withdraw(ctx, userID, withdrawAmount)

		assert.Error(t, err)
		assert.Equal(t, ErrInsufficientFunds, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestWalletService_GetBalance 測試 GetBalance 方法
func TestWalletService_GetBalance(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	initialBalance := int64(1000)
	currency := "USD"

	existingWallet := Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful get balance", func(t *testing.T) {
		mockRepo.On("GetWalletByUserID", ctx, userID).Return(existingWallet, nil)

		balance, err := service.GetBalance(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, initialBalance, balance)

		mockRepo.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockRepo.On("GetWalletByUserID", ctx, "userempty").Return(Wallet{}, ErrWalletNotFound)

		balance, err := service.GetBalance(ctx, "userempty")

		assert.Error(t, err)
		assert.Equal(t, int64(0), balance)
		assert.Equal(t, ErrWalletNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}
