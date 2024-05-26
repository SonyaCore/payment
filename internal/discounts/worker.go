package discounts

import (
	"context"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/internal/wallets"
	"payment/pkg/db"
	"strings"
)

type Worker struct {
	dataChan            chan *Seed
	DiscountService     IDiscount
	DiscountTransaction IDiscountTransaction
	WalletService       wallets.IWallet
	logger              *log.Logger
}

type Seed struct {
	ctx         context.Context
	discount    *models.Discount
	phoneNumber string
	respChan    chan *Response
}

type Response struct {
	isDone bool
	Error  error
}

// NewWorker creates a new Worker instance for charging wallet in the background.
func NewWorker(logger *log.Logger, db *db.DB, discountService IDiscount, discountTransaction IDiscountTransaction) Worker {
	return Worker{
		logger:              logger,
		DiscountService:     discountService,
		DiscountTransaction: discountTransaction,
		WalletService:       wallets.NewWallet(logger, db),
		dataChan:            make(chan *Seed),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			seed, ok := <-w.dataChan
			if !ok {
				break
			}
			response := w.Consume(seed)
			select {
			case <-seed.ctx.Done():
				return
			default:
				seed.respChan <- response
			}
		}
	}()
}

func (w *Worker) Consume(seed *Seed) *Response {
	var (
		err error
	)
	err = w.Allocation(seed.ctx, seed.discount, seed.phoneNumber)
	if err != nil {
		return &Response{
			Error: err,
		}
	}

	return &Response{
		isDone: true,
	}
}

func (w *Worker) Allocation(ctx context.Context, discount *models.Discount, phoneNumber string) error {
	var (
		err                 error
		wallet              *models.Wallet
		discountTransaction *models.DiscountTransaction
	)
	if wallet, err = w.FetchWallet(ctx, phoneNumber); err != nil {
		return err
	}

	if discountTransaction, err = w.DiscountTransaction.Add(ctx, &models.DiscountTransaction{
		DiscountID: discount.ID,
		WalletID:   wallet.ID,
		PhoneNum:   phoneNumber,
	}); err != nil {
		return err
	}

	if err = w.ChargeWallet(ctx, discount, wallet, discountTransaction); err != nil {
		return err
	}

	return nil
}

func (w *Worker) FetchWallet(ctx context.Context, phoneNumber string) (*models.Wallet, error) {
	var (
		err    error
		wallet *models.Wallet
	)
	if wallet, err = w.WalletService.GetByPhone(ctx, phoneNumber); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			if wallet, err = w.WalletService.Create(ctx, &models.Wallet{Phone: phoneNumber, Amount: 0}); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return wallet, nil
}

func (w *Worker) ChargeWallet(ctx context.Context, discount *models.Discount, wallet *models.Wallet, discountTransaction *models.DiscountTransaction) error {
	var transaction = &models.Transaction{
		WalletID:    wallet.ID,
		Amount:      discount.Amount,
		Description: discount.Description,
		Type:        models.Deposit,
	}

	if err := w.WalletService.Transaction(ctx, wallet, transaction); err != nil {
		w.logger.WithFields(log.Fields{
			"wallet_id":     wallet.ID,
			"discount_code": discount.Code,
		}).WithError(err).Error("failed to deposit wallet")

		if rollbackErr := w.DiscountTransaction.Delete(ctx, discountTransaction.ID); rollbackErr != nil {
			w.logger.WithFields(log.Fields{
				"transaction_id": discountTransaction.ID,
			}).WithError(rollbackErr).Error("failed to rollback usage")
			return rollbackErr
		}

		return err
	}

	w.logger.WithFields(log.Fields{
		"transaction_id": discountTransaction.ID,
		"amount":         discount.Amount,
		"type":           discount.Type,
	}).Info("successfully charged wallet")

	return nil
}
