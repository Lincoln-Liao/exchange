package wallet

import (
	"time"
)

type Wallet struct {
	UserID    string    // UserID is the unique identifier for the user (UUID).
	Balance   int64     // Balance is the current balance of the wallet in the smallest unit of the currency.
	Currency  string    // Currency is the type of currency the wallet holds (e.g., USD, EUR).
	CreatedAt time.Time // CreatedAt is the timestamp when the wallet was created.
	UpdatedAt time.Time // UpdatedAt is the timestamp when the wallet was last updated.
}

func NewWallet(userID, currency string) Wallet {
	return Wallet{
		UserID:    userID,
		Balance:   0,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (w *Wallet) AddBalance(amount int64) {
	w.Balance += amount
	w.UpdatedAt = time.Now()
}

func (w *Wallet) SubtractBalance(amount int64) {
	w.Balance -= amount
	w.UpdatedAt = time.Now()
}
