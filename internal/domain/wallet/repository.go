package wallet

import (
	"context"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, w Wallet) error

	GetWalletByUserID(ctx context.Context, userID string) (Wallet, error)

	UpdateWallet(ctx context.Context, w Wallet) error
}
