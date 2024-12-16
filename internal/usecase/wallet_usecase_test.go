package usecase_test

import (
	"context"
	"errors"
	"testing"

	"exchange/internal/domain/transaction"
	"exchange/internal/domain/wallet"
	"exchange/internal/usecase"
)

type mockWalletService struct {
	depositErr  error
	withdrawErr error
	balance     int64
	balanceErr  error
}

func (m *mockWalletService) CreateNewWallet(ctx context.Context, userID, currency string) (wallet.Wallet, error) {
	return wallet.Wallet{}, nil
}
func (m *mockWalletService) Deposit(ctx context.Context, userID string, amount int64) error {
	return m.depositErr
}
func (m *mockWalletService) Withdraw(ctx context.Context, userID string, amount int64) error {
	return m.withdrawErr
}
func (m *mockWalletService) GetBalance(ctx context.Context, userID string) (int64, error) {
	if m.balanceErr != nil {
		return 0, m.balanceErr
	}
	return m.balance, nil
}

type mockTransactionService struct {
	logErr            error
	getHistoryResults []transaction.Transaction
	getHistoryErr     error
}

func (m *mockTransactionService) LogTransaction(ctx context.Context, fromUserID, toUserID string, amount int64, currency string, tType transaction.TransactionType) (transaction.Transaction, error) {
	if m.logErr != nil {
		return transaction.Transaction{}, m.logErr
	}
	return transaction.Transaction{ID: "tx_mock"}, nil
}
func (m *mockTransactionService) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	if m.getHistoryErr != nil {
		return nil, m.getHistoryErr
	}
	return m.getHistoryResults, nil
}
func (m *mockTransactionService) GetTransactionByID(ctx context.Context, id string) (transaction.Transaction, error) {
	return transaction.Transaction{}, nil
}

type mockTransactionManager struct {
	doErr  error
	didRun bool
}

func (m *mockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	m.didRun = true
	if m.doErr != nil {
		return m.doErr
	}
	return fn(ctx)
}

func TestWalletUseCase_Deposit(t *testing.T) {
	wMock := &mockWalletService{}
	tMock := &mockTransactionService{}
	txMock := &mockTransactionManager{}

	uc := usecase.NewWalletUseCase(wMock, tMock, txMock)
	err := uc.Deposit(context.Background(), "user_1", 100, "TWD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !txMock.didRun {
		t.Error("expected transaction block to run")
	}

	// WalletService fail
	wMock.depositErr = wallet.ErrInvalidAmount
	err = uc.Deposit(context.Background(), "user_1", 0, "TWD")
	if err != wallet.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

	// TransactionService fail
	wMock.depositErr = nil
	tMock.logErr = errors.New("log failed")
	err = uc.Deposit(context.Background(), "user_1", 100, "TWD")
	if err == nil || err.Error() != "log failed" {
		t.Errorf("expected log failed error, got %v", err)
	}
}

func TestWalletUseCase_Withdraw(t *testing.T) {
	wMock := &mockWalletService{}
	tMock := &mockTransactionService{}
	txMock := &mockTransactionManager{}

	uc := usecase.NewWalletUseCase(wMock, tMock, txMock)
	err := uc.Withdraw(context.Background(), "user_1", 50, "TWD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wMock.withdrawErr = wallet.ErrInsufficientFunds
	err = uc.Withdraw(context.Background(), "user_1", 1000, "TWD")
	if err != wallet.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}

	// TransactionService 出錯
	wMock.withdrawErr = nil
	tMock.logErr = errors.New("log error")
	err = uc.Withdraw(context.Background(), "user_1", 50, "TWD")
	if err == nil || err.Error() != "log error" {
		t.Errorf("expected log error, got %v", err)
	}
}

func TestWalletUseCase_Transfer(t *testing.T) {
	wMock := &mockWalletService{}
	tMock := &mockTransactionService{}
	txMock := &mockTransactionManager{}
	uc := usecase.NewWalletUseCase(wMock, tMock, txMock)

	err := uc.Transfer(context.Background(), "fromUser", "toUser", 200, "TWD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// if deposit toUser fail, should rollback
	wMock.withdrawErr = nil
	wMock.depositErr = wallet.ErrInvalidAmount
	err = uc.Transfer(context.Background(), "fromUser", "toUser", 200, "TWD")
	if err != wallet.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

}

func TestWalletUseCase_GetBalance(t *testing.T) {
	wMock := &mockWalletService{balance: 500}
	tMock := &mockTransactionService{}
	txMock := &mockTransactionManager{}
	uc := usecase.NewWalletUseCase(wMock, tMock, txMock)

	bal, err := uc.GetBalance(context.Background(), "user_1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if bal != 500 {
		t.Errorf("expected balance 500, got %d", bal)
	}

	wMock.balanceErr = wallet.ErrWalletNotFound
	_, err = uc.GetBalance(context.Background(), "user_2")
	if err != wallet.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got %v", err)
	}
}

func TestWalletUseCase_GetTransactionHistory(t *testing.T) {
	wMock := &mockWalletService{}
	tMock := &mockTransactionService{
		getHistoryResults: []transaction.Transaction{
			{ID: "tx_1"},
			{ID: "tx_2"},
		},
	}
	txMock := &mockTransactionManager{}
	uc := usecase.NewWalletUseCase(wMock, tMock, txMock)

	txs, err := uc.GetTransactionHistory(context.Background(), "user_1", 10, 0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(txs))
	}

	tMock.getHistoryErr = errors.New("history error")
	_, err = uc.GetTransactionHistory(context.Background(), "user_1", 10, 0)
	if err == nil || err.Error() != "history error" {
		t.Errorf("expected 'history error', got %v", err)
	}
}
