package discounts

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/pkg/db"
)

// IDiscount defines the interface for the discount service, including methods for creating, retrieving, and checking discounts.
type IDiscount interface {
	Create(ctx context.Context, discount *models.Discount) (*models.Discount, error)
	GetByCode(ctx context.Context, code string) (*models.Discount, error)
	IsUsed(ctx context.Context, id uuid.UUID, phoneNumber string) (bool, error)
	Count(ctx context.Context, id uuid.UUID) (int64, error)
}

// NewDiscountService creates a new instance of DiscountService implementing the IDiscount interface.
func NewDiscountService(config *Config, logger *log.Logger, db *db.DB) IDiscount {
	return &DiscountService{logger, db, config}
}

// IDiscountTransaction defines the interface for the discount transaction service,
type IDiscountTransaction interface {
	List(ctx context.Context, id uuid.UUID) ([]*models.DiscountTransaction, error)
	Add(ctx context.Context, transaction *models.DiscountTransaction) (*models.DiscountTransaction, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewDiscountTransactionService creates a new instance of DiscountTransactionService implementing the IDiscountTransaction interface.
func NewDiscountTransactionService(config *Config, logger *log.Logger, db *db.DB) IDiscountTransaction {
	return &DiscountService{logger, db, config}
}

// DiscountService is the struct that implements the IDiscount and IDiscountTransaction interfaces.
type DiscountService struct {
	logger *log.Logger
	db     *db.DB
	config *Config
}

func (r *DiscountService) Create(ctx context.Context, discount *models.Discount) (*models.Discount, error) {
	if err := r.db.WithContext(ctx).Save(discount).Error; err != nil {
		r.logger.WithFields(log.Fields{
			"discount": discount,
			"error":    err,
		}).Error("failed to save discount")
		return nil, err
	}
	r.logger.WithFields(log.Fields{
		"discount": discount,
	}).Info("discount created")

	return discount, nil
}

func (r *DiscountService) GetByCode(ctx context.Context, code string) (*models.Discount, error) {
	var discount *models.Discount
	var err error
	if err := r.db.
		Model(new(models.Discount)).
		WithContext(ctx).
		Where("code = ?", code).
		First(&discount).Error; err != nil {
		r.logger.WithFields(log.Fields{
			"discount_code": code,
			"error":         err,
		}).Error("failed to find discount by Code")

		return nil, err
	}

	discount.Transactions, err = r.List(ctx, discount.ID)
	if err != nil {
		r.logger.WithFields(log.Fields{
			"discount_id": discount.ID,
			"error":       err,
		}).Error("failed to list discount transactions")
	}

	return discount, nil
}

func (r *DiscountService) IsUsed(ctx context.Context, id uuid.UUID, phoneNumber string) (bool, error) {
	var count int64 = 0
	if err := r.db.
		Model(new(models.DiscountTransaction)).
		WithContext(ctx).
		Where("discount_id = ?", id).
		Where("phone_num = ?", phoneNumber).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DiscountService) Count(ctx context.Context, id uuid.UUID) (int64, error) {
	var count int64 = 0
	if err := r.db.
		Model(new(models.DiscountTransaction)).
		WithContext(ctx).
		Where("discount_id = ?", id).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DiscountService) Add(ctx context.Context, transaction *models.DiscountTransaction) (*models.DiscountTransaction, error) {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		r.logger.WithFields(log.Fields{
			"transaction_id": transaction.ID,
			"discount_id":    transaction.DiscountID,
			"wallet_id":      transaction.WalletID,
			"error":          err,
		}).Error("failed to save discount transaction")

		return nil, err
	}
	return transaction, nil
}

func (r *DiscountService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Unscoped().WithContext(ctx).Model(new(models.DiscountTransaction)).Delete(new(models.Status), id).Error; err != nil {
		r.logger.WithFields(log.Fields{
			"transaction_id": id,
			"error":          err,
		}).Error("failed to delete discount transaction")
		return err
	}
	return nil
}

func (r *DiscountService) List(ctx context.Context, id uuid.UUID) ([]*models.DiscountTransaction, error) {
	var transactions []*models.DiscountTransaction
	if err := r.db.
		Model(new(models.DiscountTransaction)).
		WithContext(ctx).
		Where("discount_id = ?", id).
		Find(&transactions).Error; err != nil {
		r.logger.WithFields(log.Fields{
			"discount_id": id,
			"error":       err,
		}).Error("failed to list discount transactions")
		return nil, err
	}

	return transactions, nil
}
