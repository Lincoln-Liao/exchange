package wallet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"exchange/internal/domain/wallet"
)

type mockWalletRepository struct {
	wallets   map[string]wallet.Wallet
	createErr error
	updateErr error
	getErr    error
}

func newMockWalletRepository() *mockWalletRepository {
	return &mockWalletRepository{
		wallets: make(map[string]wallet.Wallet),
	}
}

func (m *mockWalletRepository) CreateWallet(ctx context.Context, w wallet.Wallet) error {
	if m.createErr != nil {
		return m.createErr
	}
	if _, exists := m.wallets[w.UserID]; exists {
		return errors.New("wallet already exists")
	}
	m.wallets[w.UserID] = w
	return nil
}

func (m *mockWalletRepository) GetWalletByUserID(ctx context.Context, userID string) (wallet.Wallet, error) {
	if m.getErr != nil {
		return wallet.Wallet{}, m.getErr
	}
	w, exists := m.wallets[userID]
	if !exists {
		return wallet.Wallet{}, wallet.ErrWalletNotFound
	}
	return w, nil
}

func (m *mockWalletRepository) UpdateWallet(ctx context.Context, w wallet.Wallet) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	_, exists := m.wallets[w.UserID]
	if !exists {
		return wallet.ErrWalletNotFound
	}
	m.wallets[w.UserID] = w
	return nil
}

func TestCreateNewWallet(t *testing.T) {
	repo := newMockWalletRepository()
	service := wallet.NewWalletService(repo)

	ctx := context.Background()
	w, err := service.CreateNewWallet(ctx, "user_1", "TWD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.UserID != "user_1" || w.Currency != "TWD" {
		t.Errorf("wallet created with wrong data: %+v", w)
	}

	_, err = service.CreateNewWallet(ctx, "user_1", "TWD")
	if err == nil {
		t.Error("expected error for duplicate wallet, got nil")
	}
}

func TestDeposit(t *testing.T) {
	repo := newMockWalletRepository()
	service := wallet.NewWalletService(repo)

	ctx := context.Background()
	w := wallet.NewWallet("user_2", "TWD")
	repo.wallets["user_2"] = w

	if err := service.Deposit(ctx, "user_2", 100); err != nil {
		t.Errorf("unexpected error on deposit: %v", err)
	}
	w2, _ := repo.GetWalletByUserID(ctx, "user_2")
	if w2.Balance != 100 {
		t.Errorf("expected balance 100, got %d", w2.Balance)
	}

	err := service.Deposit(ctx, "user_2", 0)
	if err != wallet.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

	err = service.Deposit(ctx, "nonexistent_user", 50)
	if err != wallet.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got %v", err)
	}
}

func TestWithdraw(t *testing.T) {
	repo := newMockWalletRepository()
	service := wallet.NewWalletService(repo)

	ctx := context.Background()
	w := wallet.Wallet{
		UserID:    "user_3",
		Balance:   200,
		Currency:  "TWD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.wallets["user_3"] = w

	if err := service.Withdraw(ctx, "user_3", 100); err != nil {
		t.Errorf("unexpected error on withdraw: %v", err)
	}
	w3, _ := repo.GetWalletByUserID(ctx, "user_3")
	if w3.Balance != 100 {
		t.Errorf("expected balance 100 after withdraw, got %d", w3.Balance)
	}

	err := service.Withdraw(ctx, "user_3", 200)
	if err != wallet.ErrInsufficientFunds {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}

	err = service.Withdraw(ctx, "user_3", 0)
	if err != wallet.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}

	err = service.Withdraw(ctx, "nonexistent_user", 50)
	if err != wallet.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got %v", err)
	}
}

func TestGetBalance(t *testing.T) {
	repo := newMockWalletRepository()
	service := wallet.NewWalletService(repo)

	ctx := context.Background()
	w := wallet.Wallet{
		UserID:    "user_4",
		Balance:   300,
		Currency:  "TWD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.wallets["user_4"] = w

	bal, err := service.GetBalance(ctx, "user_4")
	if err != nil {
		t.Fatalf("unexpected error on getBalance: %v", err)
	}
	if bal != 300 {
		t.Errorf("expected balance 300, got %d", bal)
	}

	_, err = service.GetBalance(ctx, "nonexistent_user")
	if err != wallet.ErrWalletNotFound {
		t.Errorf("expected ErrWalletNotFound, got %v", err)
	}
}
