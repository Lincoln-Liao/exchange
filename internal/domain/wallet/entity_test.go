package wallet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewWallet 測試 NewWallet 函數的正確性
func TestNewWallet(t *testing.T) {
	userID := "user123"
	currency := "USD"

	wallet := NewWallet(userID, currency)

	assert.Equal(t, userID, wallet.UserID, "UserID 應該正確設置")
	assert.Equal(t, int64(0), wallet.Balance, "初始餘額應該為 0")
	assert.Equal(t, currency, wallet.Currency, "Currency 應該正確設置")

	// 驗證 CreatedAt 和 UpdatedAt 是否在合理的時間範圍內
	now := time.Now()
	assert.WithinDuration(t, now, wallet.CreatedAt, time.Second, "CreatedAt 應該設置為當前時間")
	assert.WithinDuration(t, now, wallet.UpdatedAt, time.Second, "UpdatedAt 應該設置為當前時間")
}

// TestAddBalance 測試 AddBalance 方法的正確性
func TestAddBalance(t *testing.T) {
	userID := "user123"
	currency := "USD"
	initialBalance := int64(1000)
	wallet := Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Add positive amount", func(t *testing.T) {
		amountToAdd := int64(500)
		oldUpdatedAt := wallet.UpdatedAt

		wallet.AddBalance(amountToAdd)

		assert.Equal(t, initialBalance+amountToAdd, wallet.Balance, "餘額應該增加指定金額")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})

	t.Run("Add zero amount", func(t *testing.T) {
		amountToAdd := int64(0)
		oldBalance := wallet.Balance
		oldUpdatedAt := wallet.UpdatedAt

		wallet.AddBalance(amountToAdd)

		assert.Equal(t, oldBalance, wallet.Balance, "餘額應該保持不變")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})

	t.Run("Add negative amount", func(t *testing.T) {
		amountToAdd := int64(-200)
		oldBalance := wallet.Balance
		oldUpdatedAt := wallet.UpdatedAt

		wallet.AddBalance(amountToAdd)

		assert.Equal(t, oldBalance+amountToAdd, wallet.Balance, "餘額應該減少指定金額")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})
}

// TestSubtractBalance 測試 SubtractBalance 方法的正確性
func TestSubtractBalance(t *testing.T) {
	userID := "user123"
	currency := "USD"
	initialBalance := int64(1000)
	wallet := Wallet{
		UserID:    userID,
		Balance:   initialBalance,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Subtract positive amount", func(t *testing.T) {
		amountToSubtract := int64(300)
		oldUpdatedAt := wallet.UpdatedAt

		wallet.SubtractBalance(amountToSubtract)

		assert.Equal(t, initialBalance-amountToSubtract, wallet.Balance, "餘額應該減少指定金額")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})

	t.Run("Subtract zero amount", func(t *testing.T) {
		amountToSubtract := int64(0)
		oldBalance := wallet.Balance
		oldUpdatedAt := wallet.UpdatedAt

		wallet.SubtractBalance(amountToSubtract)

		assert.Equal(t, oldBalance, wallet.Balance, "餘額應該保持不變")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})

	t.Run("Subtract negative amount", func(t *testing.T) {
		amountToSubtract := int64(-100)
		oldBalance := wallet.Balance
		oldUpdatedAt := wallet.UpdatedAt

		wallet.SubtractBalance(amountToSubtract)

		assert.Equal(t, oldBalance-amountToSubtract, wallet.Balance, "餘額應該增加指定金額")
		assert.True(t, wallet.UpdatedAt.After(oldUpdatedAt), "UpdatedAt 應該更新為更晚的時間")
	})
}
