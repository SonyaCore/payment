package discounts

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"payment/api/models"
	"payment/pkg/db"
	"payment/pkg/utils"
	"time"
)

type Handler struct {
	discount    IDiscount
	transaction IDiscountTransaction
	service     Service
	logger      *log.Logger
	config      *Config
}

func NewHandler(config *Config, logger *log.Logger, db *db.DB) *Handler {
	handler := &Handler{
		discount:    NewDiscountService(config, logger, db),
		transaction: NewDiscountTransactionService(config, logger, db),
		service:     NewService(config, db, logger),
		logger:      logger,
		config:      config,
	}
	return handler
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	discountRoutes := router.PathPrefix("/discount").Subrouter()

	discountRoutes.HandleFunc("/", h.createDiscount).Methods(http.MethodPost)
	discountRoutes.HandleFunc("/usages", h.discountTransactions).Methods(http.MethodGet)
	discountRoutes.HandleFunc("/apply", h.applyDiscount).Methods(http.MethodPost)

}

// createDiscount handles the creation of a new discount code.
func (h *Handler) createDiscount(w http.ResponseWriter, r *http.Request) {
	var discount *models.Discount
	var err error

	if err = json.NewDecoder(r.Body).Decode(&discount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if discount.Type == "" {
		http.Error(w, "discount type is required", http.StatusBadRequest)
		return
	}

	switch discount.Type {
	case models.Voucher, models.Charge:
		break
	default:
		http.Error(w, "discount type is invalid", http.StatusBadRequest)
		return
	}

	discount.ExpirationTime = time.Now().Add(h.config.CreditExpiration)
	discount.CreatedAt = time.Now()
	discount.Code = utils.GenerateDiscount(h.config.CodeLength)

	if discount, err = h.discount.Create(r.Context(), discount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(discount)
	}
}

// discountTransactions handles fetching discount details by its code.
func (h *Handler) discountTransactions(w http.ResponseWriter, r *http.Request) {
	discountCode := r.URL.Query().Get("code")
	if discountCode == "" {
		http.Error(w, "discount code is required", http.StatusBadRequest)
		return
	}

	discount, err := h.discount.GetByCode(r.Context(), discountCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discount)

}

// applyDiscount applies a transaction for a given discount code and phone number.
func (h *Handler) applyDiscount(w http.ResponseWriter, r *http.Request) {
	var response = models.DiscountResponse{}
	discountCode := r.URL.Query().Get("code")
	phoneNumber := r.URL.Query().Get("phone")

	if discountCode == "" {
		http.Error(w, "discount code is required", http.StatusBadRequest)
		return
	}
	if phoneNumber == "" {
		http.Error(w, "phone number is required", http.StatusBadRequest)
		return
	}

	req := &models.DiscountApplyRequest{
		Code:     discountCode,
		PhoneNum: phoneNumber,
	}

	discount, err := h.service.Apply(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to apply discount: "+err.Error(), http.StatusBadRequest)
		return
	}

	response = models.DiscountResponse{
		Code:        discountCode,
		Description: discount.Description,
		Total:       discount.Amount,
		Type:        discount.Type,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
