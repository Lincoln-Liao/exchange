package transaction_test

import (
	"testing"
	"time"

	"exchange/internal/domain/transaction"
)

func TestNewTransaction(t *testing.T) {
	id := "tx_123"
	fromUserID := "user_from"
	toUserID := "user_to"
	amount := int64(1000)
	currency := "TWD"
	tType := transaction.TransactionTypeTransfer

	tx := transaction.NewTransaction(id, fromUserID, toUserID, amount, currency, tType)

	if tx.ID != id {
		t.Errorf("expected ID %s, got %s", id, tx.ID)
	}
	if tx.FromUserID != fromUserID {
		t.Errorf("expected FromUserID %s, got %s", fromUserID, tx.FromUserID)
	}
	if tx.ToUserID != toUserID {
		t.Errorf("expected ToUserID %s, got %s", toUserID, tx.ToUserID)
	}
	if tx.Amount != amount {
		t.Errorf("expected Amount %d, got %d", amount, tx.Amount)
	}
	if tx.Currency != currency {
		t.Errorf("expected Currency %s, got %s", currency, tx.Currency)
	}
	if tx.Type != tType {
		t.Errorf("expected Type %s, got %s", tType, tx.Type)
	}

	if tx.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestTransactionTypeConstants(t *testing.T) {
	if transaction.TransactionTypeDeposit != "DEPOSIT" {
		t.Errorf("expected TransactionTypeDeposit to be 'DEPOSIT', got '%s'", transaction.TransactionTypeDeposit)
	}
	if transaction.TransactionTypeWithdraw != "WITHDRAW" {
		t.Errorf("expected TransactionTypeWithdraw to be 'WITHDRAW', got '%s'", transaction.TransactionTypeWithdraw)
	}
	if transaction.TransactionTypeTransfer != "TRANSFER" {
		t.Errorf("expected TransactionTypeTransfer to be 'TRANSFER', got '%s'", transaction.TransactionTypeTransfer)
	}
}

func TestNewTransactionInvalidAmount(t *testing.T) {
	tx := transaction.NewTransaction("tx_abc", "user_from", "user_to", 0, "TWD", transaction.TransactionTypeWithdraw)
	if tx.Amount != 0 {
		t.Errorf("expected Amount 0, got %d", tx.Amount)
	}

	tx = transaction.NewTransaction("tx_xyz", "user_from", "user_to", -100, "TWD", transaction.TransactionTypeDeposit)
	if tx.Amount != -100 {
		t.Errorf("expected Amount -100, got %d", tx.Amount)
	}
}

func TestCreatedAtIsSet(t *testing.T) {
	tx := transaction.NewTransaction("tx_time", "uf", "ut", 500, "TWD", transaction.TransactionTypeDeposit)
	if tx.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	now := time.Now()
	if tx.CreatedAt.After(now) {
		t.Error("expected CreatedAt not to be in the future")
	}
}
