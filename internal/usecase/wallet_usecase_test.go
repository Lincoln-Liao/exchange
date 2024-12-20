// wallet_usecase_test.go
package usecase

import (
	"context"
	"testing"

	"exchange/internal/domain/transaction"
	"exchange/internal/domain/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateNewWallet(ctx context.Context, userID, currency string) (wallet.Wallet, error) {
	args := m.Called(ctx, userID, currency)
	return args.Get(0).(wallet.Wallet), args.Error(1)
}

func (m *MockWalletService) Deposit(ctx context.Context, userID string, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockWalletService) Withdraw(ctx context.Context, userID string, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockWalletService) GetBalance(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

type MockTransactionManager struct {
	mock.Mock
	DoFn func(ctx context.Context, fn func(ctx context.Context) error) error
}

func (m *MockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.DoFn != nil {
		return m.DoFn(ctx, fn)
	}
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestWalletUseCase_Deposit(t *testing.T) {
	mockWalletService := new(MockWalletService)
	mockTransactionService := new(MockTransactionService)
	mockTxManager := new(MockTransactionManager)

	useCase := NewWalletUseCase(mockWalletService, mockTransactionService, mockTxManager)

	ctx := context.Background()
	userID := "user1"
	amount := int64(1000)
	currency := "USD"

	t.Run("successful deposit", func(t *testing.T) {
		mockTxManager.DoFn = func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}

		mockWalletService.On("Deposit", ctx, userID, amount).Return(nil)

		expectedTx := transaction.Transaction{
			ID:         "tx123",
			FromUserID: "",
			ToUserID:   userID,
			Amount:     amount,
			Currency:   currency,
			Type:       transaction.TransactionTypeDeposit,
		}
		mockTransactionService.On("LogTransaction", ctx, "", userID, amount, currency, transaction.TransactionTypeDeposit).Return(expectedTx, nil)

		err := useCase.Deposit(ctx, userID, amount, currency)

		assert.NoError(t, err)
		mockWalletService.AssertExpectations(t)
		mockTransactionService.AssertExpectations(t)
	})

	t.Run("invalid amount", func(t *testing.T) {
		invalidAmount := int64(-100)

		mockTxManager.DoFn = func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}

		mockWalletService.On("Deposit", ctx, userID, invalidAmount).Return(wallet.ErrInvalidAmount)

		mockTransactionService.On("LogTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(transaction.Transaction{}, nil).Maybe()

		err := useCase.Deposit(ctx, userID, invalidAmount, currency)

		assert.ErrorIs(t, err, wallet.ErrInvalidAmount)
		mockWalletService.AssertExpectations(t)
		mockTransactionService.AssertExpectations(t)
	})
}

func TestWalletUseCase_Withdraw(t *testing.T) {
	mockWalletService := new(MockWalletService)
	mockTransactionService := new(MockTransactionService)
	mockTxManager := new(MockTransactionManager)

	useCase := NewWalletUseCase(mockWalletService, mockTransactionService, mockTxManager)

	ctx := context.Background()
	userID := "user1"
	amount := int64(500)
	currency := "USD"

	t.Run("successful withdraw", func(t *testing.T) {
		mockTxManager.DoFn = func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}

		mockWalletService.On("Withdraw", ctx, userID, amount).Return(nil)

		expectedTx := transaction.Transaction{
			ID:         "tx124",
			FromUserID: userID,
			ToUserID:   "",
			Amount:     amount,
			Currency:   currency,
			Type:       transaction.TransactionTypeWithdraw,
		}
		mockTransactionService.On("LogTransaction", ctx, userID, "", amount, currency, transaction.TransactionTypeWithdraw).Return(expectedTx, nil)

		err := useCase.Withdraw(ctx, userID, amount, currency)

		assert.NoError(t, err)
		mockTxManager.AssertExpectations(t)
		mockWalletService.AssertExpectations(t)
		mockTransactionService.AssertExpectations(t)
	})

	t.Run("invalid amount", func(t *testing.T) {
		invalidAmount := int64(-200)

		mockTxManager.DoFn = func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}
		mockWalletService.On("Withdraw", ctx, userID, invalidAmount).Return(wallet.ErrInvalidAmount)

		mockTransactionService.On("LogTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(transaction.Transaction{}, nil).Maybe()

		err := useCase.Withdraw(ctx, userID, invalidAmount, currency)

		assert.ErrorIs(t, err, wallet.ErrInvalidAmount)
		mockTxManager.AssertExpectations(t)
		mockWalletService.AssertExpectations(t)
		mockTransactionService.AssertExpectations(t)
	})
}

func TestWalletUseCase_Transfer(t *testing.T) {
	mockWalletService := new(MockWalletService)
	mockTransactionService := new(MockTransactionService)
	mockTxManager := new(MockTransactionManager)

	useCase := NewWalletUseCase(mockWalletService, mockTransactionService, mockTxManager)

	ctx := context.Background()
	fromUserID := "user1"
	toUserID := "user2"
	amount := int64(300)
	currency := "USD"

	t.Run("successful transfer", func(t *testing.T) {
		mockTxManager.DoFn = func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}
		mockWalletService.On("Withdraw", ctx, fromUserID, amount).Return(nil)
		mockWalletService.On("Deposit", ctx, toUserID, amount).Return(nil)

		expectedTx := transaction.Transaction{
			ID:         "tx125",
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Currency:   currency,
			Type:       transaction.TransactionTypeTransfer,
		}
		mockTransactionService.On("LogTransaction", ctx, fromUserID, toUserID, amount, currency, transaction.TransactionTypeTransfer).Return(expectedTx, nil)

		err := useCase.Transfer(ctx, fromUserID, toUserID, amount, currency)

		assert.NoError(t, err)
		mockTxManager.AssertExpectations(t)
		mockWalletService.AssertExpectations(t)
		mockTransactionService.AssertExpectations(t)
	})
}

func TestWalletUseCase_GetBalance(t *testing.T) {
	mockWalletService := new(MockWalletService)
	mockTransactionService := new(MockTransactionService)
	mockTxManager := new(MockTransactionManager)

	useCase := NewWalletUseCase(mockWalletService, mockTransactionService, mockTxManager)

	ctx := context.Background()
	userID := "user1"
	expectedBalance := int64(1500)

	t.Run("successful get balance", func(t *testing.T) {
		mockWalletService.On("GetBalance", ctx, userID).Return(expectedBalance, nil)

		balance, err := useCase.GetBalance(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedBalance, balance)
		mockWalletService.AssertExpectations(t)
	})

	t.Run("wallet service get balance error", func(t *testing.T) {
		getBalanceErr := wallet.ErrWalletNotFound

		mockWalletService.On("GetBalance", ctx, "userIDEmpty").Return(int64(0), getBalanceErr)

		balance, err := useCase.GetBalance(ctx, "userIDEmpty")

		assert.ErrorIs(t, err, getBalanceErr)
		assert.Equal(t, int64(0), balance)
		mockWalletService.AssertExpectations(t)
	})
}
