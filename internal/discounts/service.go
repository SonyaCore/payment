package discounts

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/internal/wallets"
	"payment/pkg/db"
	"time"
)

// Service aggregates various services and components required for transactions, charging wallets operations.
type Service struct {
	discountService     IDiscount
	discountTransaction IDiscountTransaction
	walletService       wallets.IWallet
	worker              Worker
	configs             *Config
	logger              *log.Logger
}

// NewService initializes and returns a new Service instance.
func NewService(config *Config, db *db.DB, log *log.Logger) Service {
	discountService := NewDiscountService(config, log, db)
	discountTransaction := NewDiscountTransactionService(config, log, db)
	service := Service{
		discountService:     discountService,
		discountTransaction: discountTransaction,
		walletService:       wallets.NewWallet(log, db),
		worker:              NewWorker(log, db, discountService, discountTransaction),
		configs:             config,
		logger:              log,
	}

	// Start the worker to handle background tasks
	service.worker.Start()

	return service
}

func (s *Service) Create(ctx context.Context, discount *models.Discount) (*models.Discount, error) {
	s.logger.WithFields(log.Fields{
		"discount": discount,
	}).Info("discount created successfully")

	if tx, err := s.discountService.Create(ctx, discount); err != nil {
		return nil, err
	} else {
		return tx, nil
	}
}

func (s *Service) Apply(ctx context.Context, req *models.DiscountApplyRequest) (*models.Discount, error) {
	var (
		err      error
		discount *models.Discount
	)

	if discount, err = s.discountService.GetByCode(ctx, req.Code); err != nil {
		return nil, err
	}

	if err = s.IsUsed(ctx, discount, req.PhoneNum); err != nil {
		return nil, err
	}

	if err = s.IsExpired(ctx, discount); err != nil {
		return nil, err
	}

	timeout := time.Tick(s.configs.CreditExpiration)
	workerResp := make(chan *Response)

	s.worker.dataChan <- &Seed{
		ctx:         ctx,
		discount:    discount,
		phoneNumber: req.PhoneNum,
		respChan:    workerResp,
	}

	select {
	case resp := <-workerResp:
		if resp.isDone {
			return discount, nil
		}
		return nil, resp.Error
	case <-ctx.Done():
		return nil, err
	case <-timeout:
		return nil, err
	}
}

func (s *Service) IsUsed(ctx context.Context, discount *models.Discount, phoneNumber string) error {
	var (
		used       bool
		usageCount int64
		err        error
	)

	if usageCount, err = s.discountService.Count(ctx, discount.ID); err != nil {
		return err
	}

	if usageCount >= discount.UsageLimit {
		return errors.New("usage limit exceed")
	}

	if used, err = s.discountService.IsUsed(ctx, discount.ID, phoneNumber); err != nil {
		return err
	}
	if used {
		return errors.New("discount code used before")
	}
	return nil
}

func (s *Service) IsExpired(_ context.Context, discount *models.Discount) error {
	if discount.ExpirationTime.Before(time.Now()) {
		return errors.New("discount expired")
	}
	if discount.CreatedAt.After(time.Now()) {
		return errors.New("invalid discount")
	}
	return nil
}
