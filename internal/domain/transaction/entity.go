package transaction

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID         string          // Unique transaction identifier
	FromUserID string          // Source user ID
	ToUserID   string          // Target user ID
	Amount     int64           // Transaction amount, expressed as an integer in the smallest currency unit
	Currency   string          // Currency code (e.g., "USD", "TWD")
	Type       TransactionType // Transaction type (DEPOSIT, WITHDRAW, TRANSFER)
	CreatedAt  time.Time       // Transaction creation time
}

func NewTransaction(id, fromUserID, toUserID string, amount int64, currency string, tType TransactionType) (Transaction, error) {
	if id == "" {
		return Transaction{}, ErrInvalidTransactionID
	}
	if amount <= 0 {
		return Transaction{}, ErrInvalidTransactionAmount
	}
	return Transaction{
		ID:         id,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Currency:   currency,
		Type:       tType,
		CreatedAt:  time.Now(),
	}, nil
}
