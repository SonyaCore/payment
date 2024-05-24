package wallets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"payment/api/models"
	"payment/internal/transactions"
	"payment/pkg/db"
)

// RegisterRoutes registers the routes for wallet-related operations with the provided router.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	walletRoutes := router.PathPrefix("/wallet").Subrouter()

	walletRoutes.HandleFunc("/register", h.createWalletHandler).Methods(http.MethodPost)
	walletRoutes.HandleFunc("/{phoneNumber}", h.transactionHandler).Methods(http.MethodPut)
	walletRoutes.HandleFunc("/{phoneNumber}", h.deleteWalletHandler).Methods(http.MethodDelete)
	walletRoutes.HandleFunc("/{phoneNumber}", h.returnByPhoneNumber).Methods(http.MethodGet)
}

// Handler is a struct that holds the services and logger needed for handling wallet and transaction-related requests.
type Handler struct {
	WalletService      IWallet
	TransactionService transactions.ITransaction
	Logger             *logrus.Logger
}

// NewHandler initializes a new Handler with the provided database connection and logger.
// It sets up the WalletService and TransactionService with their respective dependencies.
func NewHandler(db *db.DB, logger *logrus.Logger) *Handler {
	handler := &Handler{
		Logger:             logger,
		TransactionService: transactions.NewTransactionsService(logger, db),
		WalletService:      NewWallet(logger, db),
	}
	return handler
}

// createWalletHandler handles the creation of a new wallet.
func (h *Handler) createWalletHandler(w http.ResponseWriter, r *http.Request) {
	var wallet *models.Wallet
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if model, _ := h.WalletService.GetByPhone(ctx, wallet.Phone); model != nil {
		http.Error(w, "wallet already exist", http.StatusBadRequest)
		return
	}

	wallet, err = h.WalletService.Create(ctx, wallet)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// deleteWalletHandler handles the deletion of a wallet by phone number.
func (h *Handler) deleteWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var phoneNumber = vars["phoneNumber"]

	ctx := context.Background()
	wallet, err := h.WalletService.Delete(ctx, phoneNumber)

	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// transactionHandler handles wallet transactions (withdrawals and deposits).
func (h *Handler) transactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction *models.NewTransaction
	vars := mux.Vars(r)
	var phoneNumber = vars["phoneNumber"]

	ctx := context.Background()
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if phoneNumber == "" {
		http.Error(w, "wallet phone required", http.StatusBadRequest)
		return
	}

	transaction.Phone = phoneNumber
	wallet, err := h.WalletService.GetByPhone(ctx, phoneNumber)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a transaction based on the request data
	var tx = &models.Transaction{
		WalletID:    wallet.ID,
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Description: transaction.Description,
	}

	switch transaction.Type {
	case models.Withdrawal:
		transaction.Type = models.Withdrawal
		if wallet.Amount < transaction.Amount {
			h.Logger.WithFields(logrus.Fields{
				"type":            "transaction",
				"wallet_id":       wallet.ID,
				"current_amount":  wallet.Amount,
				"required_amount": transaction.Amount,
			}).Error("Insufficient funds")

			http.Error(w, fmt.Sprintf("insufficient funds: wallet %s balance is %d, but %d is required", wallet.Phone, wallet.Amount, transaction.Amount), http.StatusBadRequest)
			return
		}
	case models.Deposit:
		transaction.Type = models.Deposit
	default:
		http.Error(w, "unknown transaction type", http.StatusBadRequest)
		return
	}

	err = h.WalletService.Transaction(ctx, wallet, tx)
	if err != nil {
		h.Logger.WithFields(logrus.Fields{
			"type":        "transaction",
			"wallet_id":   wallet.ID,
			"transaction": tx,
			"error":       err,
		}).Error("Transaction failed")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// returnByPhoneNumber handles retrieving a wallet by phone number.
func (h *Handler) returnByPhoneNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phoneNumber := vars["phoneNumber"]

	ctx := context.Background()
	wallet, err := h.WalletService.GetByPhone(ctx, phoneNumber)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(wallet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}