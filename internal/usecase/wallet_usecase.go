package usecase

import (
	"context"
	"exchange/internal/domain/transaction"
	"exchange/internal/domain/wallet"
)

type WalletServiceInterface interface {
	CreateNewWallet(ctx context.Context, userID, currency string) (wallet.Wallet, error)
	Deposit(ctx context.Context, userID string, amount int64) error
	Withdraw(ctx context.Context, userID string, amount int64) error
	GetBalance(ctx context.Context, userID string) (int64, error)
}

type TransactionServiceInterface interface {
	LogTransaction(ctx context.Context, fromUserID, toUserID string, amount int64, currency string, tType transaction.TransactionType) (transaction.Transaction, error)
	GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error)
	GetTransactionByID(ctx context.Context, id string) (transaction.Transaction, error)
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type WalletUseCase struct {
	walletService      WalletServiceInterface
	transactionService TransactionServiceInterface
	txManager          TransactionManager
}

func NewWalletUseCase(
	wService WalletServiceInterface,
	tService TransactionServiceInterface,
	txManager TransactionManager,
) *WalletUseCase {
	return &WalletUseCase{
		walletService:      wService,
		transactionService: tService,
		txManager:          txManager,
	}
}

func (uc *WalletUseCase) Deposit(ctx context.Context, userID string, amount int64, currency string) error {
	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.walletService.Deposit(ctx, userID, amount); err != nil {
			return err
		}
		_, err := uc.transactionService.LogTransaction(ctx, "", userID, amount, currency, transaction.TransactionTypeDeposit)
		return err
	})
}

func (uc *WalletUseCase) Withdraw(ctx context.Context, userID string, amount int64, currency string) error {
	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.walletService.Withdraw(ctx, userID, amount); err != nil {
			return err
		}
		_, err := uc.transactionService.LogTransaction(ctx, userID, "", amount, currency, transaction.TransactionTypeWithdraw)
		return err
	})
}

func (uc *WalletUseCase) Transfer(ctx context.Context, fromUserID, toUserID string, amount int64, currency string) error {
	return uc.txManager.Do(ctx, func(ctx context.Context) error {
		if err := uc.walletService.Withdraw(ctx, fromUserID, amount); err != nil {
			return err
		}

		if err := uc.walletService.Deposit(ctx, toUserID, amount); err != nil {
			return err
		}

		_, err := uc.transactionService.LogTransaction(ctx, fromUserID, toUserID, amount, currency, transaction.TransactionTypeTransfer)
		return err
	})
}

func (uc *WalletUseCase) GetBalance(ctx context.Context, userID string) (int64, error) {
	return uc.walletService.GetBalance(ctx, userID)
}

func (uc *WalletUseCase) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]transaction.Transaction, error) {
	return uc.transactionService.GetTransactionHistory(ctx, userID, limit, offset)
}
