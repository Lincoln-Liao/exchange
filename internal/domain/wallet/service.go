package wallet

import (
	"context"
)

type WalletServiceInterface interface {
	CreateNewWallet(ctx context.Context, userID, currency string) (Wallet, error)
	Deposit(ctx context.Context, userID string, amount int64) error
	Withdraw(ctx context.Context, userID string, amount int64) error
	GetBalance(ctx context.Context, userID string) (int64, error)
}

type WalletService struct {
	repository WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{
		repository: repo,
	}
}

func (s *WalletService) CreateNewWallet(ctx context.Context, userID, currency string) (Wallet, error) {
	w := NewWallet(userID, currency)
	err := s.repository.CreateWallet(ctx, w)
	if err != nil {
		return Wallet{}, err
	}
	return w, nil
}

func (s *WalletService) Deposit(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w, err := s.repository.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == ErrWalletNotFound {
			return ErrWalletNotFound
		}
		return err
	}

	w.AddBalance(amount)

	if err := s.repository.UpdateWallet(ctx, w); err != nil {
		return err
	}

	return nil
}

func (s *WalletService) Withdraw(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w, err := s.repository.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == ErrWalletNotFound {
			return ErrWalletNotFound
		}
		return err
	}

	if w.Balance < amount {
		return ErrInsufficientFunds
	}

	w.SubtractBalance(amount)

	if err := s.repository.UpdateWallet(ctx, w); err != nil {
		return err
	}

	return nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID string) (int64, error) {
	w, err := s.repository.GetWalletByUserID(ctx, userID)
	if err != nil {
		if err == ErrWalletNotFound {
			return 0, ErrWalletNotFound
		}
		return 0, ErrDatabaseFailure
	}
	return w.Balance, nil
}
