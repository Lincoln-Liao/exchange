package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"exchange/internal/domain/transaction"
	"exchange/internal/domain/wallet"
	"exchange/internal/usecase"
)

type Handler struct {
	WalletUC *usecase.WalletUseCase
}

func NewHandler(walletUC *usecase.WalletUseCase) *Handler {
	return &Handler{
		WalletUC: walletUC,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/wallet/deposit", h.depositHandler)
	mux.HandleFunc("/wallet/withdraw", h.withdrawHandler)
	mux.HandleFunc("/wallet/transfer", h.transferHandler)
	// GET /wallet/{user_id}/balance
	// GET /wallet/{user_id}/transactions?limit=10&offset=0
	mux.HandleFunc("/wallet/", h.userWalletHandler)
}

func (h *Handler) depositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.WalletUC.Deposit(ctx, req.UserID, req.Amount, req.Currency); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *Handler) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.WalletUC.Withdraw(ctx, req.UserID, req.Amount, req.Currency); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *Handler) transferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.WalletUC.Transfer(ctx, req.FromUserID, req.ToUserID, req.Amount, req.Currency); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *Handler) userWalletHandler(w http.ResponseWriter, r *http.Request) {
	// GET /wallet/{user_id}/balance
	// GET /wallet/{user_id}/transactions?limit=10&offset=0
	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/wallet/"), "/")
	if len(segments) == 0 {
		http.Error(w, "user_id not provided", http.StatusBadRequest)
		return
	}

	userID := segments[0]

	if len(segments) == 2 && segments[1] == "balance" && r.Method == http.MethodGet {
		h.getBalanceHandler(w, r, userID)
		return
	}

	if len(segments) == 2 && segments[1] == "transactions" && r.Method == http.MethodGet {
		h.getTransactionsHandler(w, r, userID)
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func (h *Handler) getBalanceHandler(w http.ResponseWriter, r *http.Request, userID string) {
	ctx := r.Context()
	balance, err := h.WalletUC.GetBalance(ctx, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := BalanceResponse{
		UserID:  userID,
		Balance: balance,
	}
	writeJSON(w, resp)
}

func (h *Handler) getTransactionsHandler(w http.ResponseWriter, r *http.Request, userID string) {
	ctx := r.Context()
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit value", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset value", http.StatusBadRequest)
			return
		}
	}

	txs, err := h.WalletUC.GetTransactionHistory(ctx, userID, limit, offset)
	if err != nil {
		handleError(w, err)
		return
	}

	resp := make([]TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		resp = append(resp, TransactionResponse{
			ID:         tx.ID,
			FromUserID: tx.FromUserID,
			ToUserID:   tx.ToUserID,
			Amount:     tx.Amount,
			Currency:   tx.Currency,
			Type:       string(tx.Type),
			CreatedAt:  tx.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	writeJSON(w, resp)
}

func handleError(w http.ResponseWriter, err error) {
	log.Println("error:", err)
	switch err {
	case wallet.ErrWalletNotFound:
		http.Error(w, "wallet not found", http.StatusNotFound)
	case wallet.ErrInsufficientFunds:
		http.Error(w, "insufficient funds", http.StatusBadRequest)
	case wallet.ErrInvalidAmount:
		http.Error(w, "invalid amount", http.StatusBadRequest)
	case transaction.ErrInvalidTransactionAmount:
		http.Error(w, "invalid transaction amount", http.StatusBadRequest)
	case transaction.ErrInvalidTransactionType:
		http.Error(w, "invalid transaction type", http.StatusBadRequest)
	case transaction.ErrTransactionNotFound:
		http.Error(w, "transaction not found", http.StatusNotFound)
	case transaction.ErrInvalidUserID:
		http.Error(w, "invalid user id", http.StatusBadRequest)
	case transaction.ErrInvalidTransactionID:
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
