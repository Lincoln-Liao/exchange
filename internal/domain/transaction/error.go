package transaction

import "errors"

var (
	ErrInvalidTransactionAmount = errors.New("invalid transaction amount")
	ErrInvalidTransactionType   = errors.New("invalid transaction type")
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrInvalidUserID            = errors.New("invalid user ID")
	ErrInvalidTransactionID     = errors.New("invalid transaction ID")
	ErrDatabaseFailure          = errors.New("database failure")
)
