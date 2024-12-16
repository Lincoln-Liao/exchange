package transaction_test

import (
	"context"
	"testing"
	"time"

	"exchange/internal/domain/transaction"
)

type mockTransactionRepository struct {
	transactions    map[string]transaction.Transaction
	createErr       error
	getByIDErr      error
	listByUserIDErr error
}

func newMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{
		transactions: make(map[string]transaction.Transaction),
	}
}

func (m *mockTransactionRepository) CreateTransaction(ctx context.Context, tx transaction.Transaction) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.transactions[tx.ID] = tx
	return nil
}

func (m *mockTransactionRepository) GetTransactionByID(ctx context.Context, id string) (transaction.Transaction, error) {
	if m.getByIDErr != nil {
		return transaction.Transaction{}, m.getByIDErr
	}
	tx, exists := m.transactions[id]
	if !exists {
		return transaction.Transaction{}, transaction.ErrTransactionNotFound
	}
	return tx, nil
}

func (m *mockTransactionRepository) ListTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	if m.listByUserIDErr != nil {
		return nil, m.listByUserIDErr
	}
	var result []transaction.Transaction
	for _, tx := range m.transactions {
		if tx.FromUserID == userID || tx.ToUserID == userID {
			result = append(result, tx)
		}
	}
	return result, nil
}

func TestLogTransaction(t *testing.T) {
	repo := newMockTransactionRepository()
	service := transaction.NewTransactionService(repo)

	ctx := context.Background()

	tx, err := service.LogTransaction(ctx, "user_from", "user_to", 100, "TWD", transaction.TransactionTypeTransfer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Amount != 100 || tx.FromUserID != "user_from" || tx.ToUserID != "user_to" || tx.Type != transaction.TransactionTypeTransfer {
		t.Errorf("transaction fields not set correctly: %+v", tx)
	}

	_, err = service.LogTransaction(ctx, "user_from", "user_to", 0, "TWD", transaction.TransactionTypeWithdraw)
	if err != transaction.ErrInvalidTransactionAmount {
		t.Errorf("expected ErrInvalidTransactionAmount, got %v", err)
	}

	_, err = service.LogTransaction(ctx, "user_from", "user_to", 100, "TWD", "INVALID_TYPE")
	if err != transaction.ErrInvalidTransactionType {
		t.Errorf("expected ErrInvalidTransactionType, got %v", err)
	}
}

func TestGetTransactionHistory(t *testing.T) {
	repo := newMockTransactionRepository()
	service := transaction.NewTransactionService(repo)

	ctx := context.Background()

	tx1 := transaction.Transaction{
		ID:         "tx_1",
		FromUserID: "userA",
		ToUserID:   "userB",
		Amount:     500,
		Currency:   "TWD",
		Type:       transaction.TransactionTypeTransfer,
		CreatedAt:  time.Now(),
	}
	tx2 := transaction.Transaction{
		ID:         "tx_2",
		FromUserID: "userB",
		ToUserID:   "userA",
		Amount:     300,
		Currency:   "TWD",
		Type:       transaction.TransactionTypeTransfer,
		CreatedAt:  time.Now(),
	}
	repo.transactions["tx_1"] = tx1
	repo.transactions["tx_2"] = tx2

	_, err := service.GetTransactionHistory(ctx, "", 10, 0)
	if err != transaction.ErrInvalidUserID {
		t.Errorf("expected ErrInvalidUserID, got %v", err)
	}

	txs, err := service.GetTransactionHistory(ctx, "userA", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(txs))
	}

	txs, err = service.GetTransactionHistory(ctx, "userC", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(txs) != 0 {
		t.Errorf("expected 0 transactions for userC, got %d", len(txs))
	}
}

func TestGetTransactionByID(t *testing.T) {
	repo := newMockTransactionRepository()
	service := transaction.NewTransactionService(repo)

	ctx := context.Background()
	txID := "tx_abc"

	_, err := service.GetTransactionByID(ctx, txID)
	if err != transaction.ErrTransactionNotFound {
		t.Errorf("expected ErrTransactionNotFound, got %v", err)
	}

	_, err = service.GetTransactionByID(ctx, "")
	if err != transaction.ErrInvalidTransactionID {
		t.Errorf("expected ErrInvalidTransactionID, got %v", err)
	}

	tx := transaction.Transaction{
		ID:         txID,
		FromUserID: "userX",
		ToUserID:   "userY",
		Amount:     1000,
		Currency:   "USD",
		Type:       transaction.TransactionTypeDeposit,
		CreatedAt:  time.Now(),
	}
	repo.transactions[txID] = tx

	gotTx, err := service.GetTransactionByID(ctx, txID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotTx.ID != txID || gotTx.Amount != 1000 || gotTx.Type != transaction.TransactionTypeDeposit {
		t.Errorf("got transaction mismatch, expected %+v, got %+v", tx, gotTx)
	}
}
