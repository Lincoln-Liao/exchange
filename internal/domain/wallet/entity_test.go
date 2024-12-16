package wallet_test

import (
	"testing"
	"time"

	"exchange/internal/domain/wallet"
)

func TestNewWallet(t *testing.T) {
	userID := "user_123"
	currency := "TWD"
	w := wallet.NewWallet(userID, currency)

	if w.UserID != userID {
		t.Errorf("expected UserID: %s, got %s", userID, w.UserID)
	}
	if w.Currency != currency {
		t.Errorf("expected Currency: %s, got %s", currency, w.Currency)
	}
	if w.Balance != 0 {
		t.Errorf("expected initial Balance: 0, got %d", w.Balance)
	}

	if w.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if w.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestAddBalance(t *testing.T) {
	w := wallet.Wallet{
		UserID:    "user_123",
		Balance:   100,
		Currency:  "TWD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.AddBalance(50)
	if w.Balance != 150 {
		t.Errorf("expected Balance: 150, got %d", w.Balance)
	}
	if w.UpdatedAt.Before(w.CreatedAt) {
		t.Error("expected UpdatedAt to be updated after AddBalance")
	}
}

func TestSubtractBalance(t *testing.T) {
	w := wallet.Wallet{
		UserID:    "user_123",
		Balance:   100,
		Currency:  "TWD",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.SubtractBalance(40)
	if w.Balance != 60 {
		t.Errorf("expected Balance: 60, got %d", w.Balance)
	}
	if w.UpdatedAt.Before(w.CreatedAt) {
		t.Error("expected UpdatedAt to be updated after SubtractBalance")
	}
}
