package transaction

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		fromUserID    string
		toUserID      string
		amount        int64
		currency      string
		tType         TransactionType
		expectError   bool
		expectedError error
	}{
		{
			name:       "valid deposit transaction",
			id:         "tx123",
			fromUserID: "",
			toUserID:   "user1",
			amount:     1000,
			currency:   "USD",
			tType:      TransactionTypeDeposit,
		},
		{
			name:       "valid withdraw transaction",
			id:         "tx124",
			fromUserID: "user1",
			toUserID:   "",
			amount:     500,
			currency:   "USD",
			tType:      TransactionTypeWithdraw,
		},
		{
			name:       "valid transfer transaction",
			id:         "tx125",
			fromUserID: "user1",
			toUserID:   "user2",
			amount:     300,
			currency:   "USD",
			tType:      TransactionTypeTransfer,
		},
		{
			name:          "invalid amount (negative)",
			id:            "tx126",
			fromUserID:    "user1",
			toUserID:      "user2",
			amount:        -100,
			currency:      "USD",
			tType:         TransactionTypeTransfer,
			expectError:   true,
			expectedError: ErrInvalidTransactionAmount,
		},
		{
			name:          "invalid amount (zero)",
			id:            "tx127",
			fromUserID:    "user1",
			toUserID:      "user2",
			amount:        0,
			currency:      "USD",
			tType:         TransactionTypeTransfer,
			expectError:   true,
			expectedError: ErrInvalidTransactionAmount,
		},
		{
			name:          "missing ID",
			id:            "",
			fromUserID:    "user1",
			toUserID:      "user2",
			amount:        100,
			currency:      "USD",
			tType:         TransactionTypeTransfer,
			expectError:   true,
			expectedError: ErrInvalidTransactionID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := NewTransaction(tt.id, tt.fromUserID, tt.toUserID, tt.amount, tt.currency, tt.tType)

			if tt.expectError {
				assert.Error(t, err, "Expected an error but got none")
				assert.EqualError(t, err, tt.expectedError.Error(), "Expected error message to match")
				// 當有錯誤時，返回的 Transaction 應該是空的
				assert.Equal(t, Transaction{}, tx, "Expected empty Transaction struct on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.Equal(t, tt.id, tx.ID, "Transaction ID should match")
				assert.Equal(t, tt.fromUserID, tx.FromUserID, "FromUserID should match")
				assert.Equal(t, tt.toUserID, tx.ToUserID, "ToUserID should match")
				assert.Equal(t, tt.amount, tx.Amount, "Amount should match")
				assert.Equal(t, tt.currency, tx.Currency, "Currency should match")
				assert.Equal(t, tt.tType, tx.Type, "TransactionType should match")

				// 驗證 CreatedAt 是否在合理的時間範圍內
				now := time.Now()
				assert.WithinDuration(t, now, tx.CreatedAt, time.Second, "CreatedAt should be within 1 second of now")
			}
		})
	}
}

func TestTransactionTypes(t *testing.T) {
	assert.Equal(t, "DEPOSIT", string(TransactionTypeDeposit), "TransactionTypeDeposit should be 'DEPOSIT'")
	assert.Equal(t, "WITHDRAW", string(TransactionTypeWithdraw), "TransactionTypeWithdraw should be 'WITHDRAW'")
	assert.Equal(t, "TRANSFER", string(TransactionTypeTransfer), "TransactionTypeTransfer should be 'TRANSFER'")
}
