package http

type DepositRequest struct {
	UserID   string `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type WithdrawRequest struct {
	UserID   string `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type TransferRequest struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
}

type BalanceResponse struct {
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}

type TransactionResponse struct {
	ID         string `json:"id"`
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
}
