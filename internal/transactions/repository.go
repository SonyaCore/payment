package transactions

import (
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/pkg/db"
)

type ITransaction interface {
	Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	CanChange(ctx context.Context, id uuid.UUID) bool
	ChangeStatus(ctx context.Context, id uuid.UUID, status models.Status) error
	List(ctx context.Context, id uuid.UUID) ([]*models.Transaction, error)
}

type Service struct {
	logger *log.Logger
	db     *db.DB
}

func NewTransactionsService(logger *log.Logger, db *db.DB) ITransaction {
	return &Service{
		logger: logger,
		db:     db,
	}
}

func (s *Service) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	if transaction.WalletID == uuid.Nil {
		return nil, errors.New("wallet ID is required")
	}
	if err := s.db.WithContext(ctx).Create(transaction).Error; err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return transaction, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	if err := s.db.WithContext(ctx).First(transaction, id).Error; err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return transaction, nil
}

func (s *Service) ChangeStatus(_ context.Context, id uuid.UUID, status models.Status) error {
	if err := s.db.
		Model(new(models.Transaction)).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		s.logger.Error(err)
		return err
	}
	return nil
}

func (s *Service) List(ctx context.Context, id uuid.UUID) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)
	if err := s.db.
		Model(new(models.Transaction)).
		WithContext(ctx).
		Where("wallet_id = ?", id).
		Find(&transactions).Error; err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return transactions, nil
}

func (s *Service) CanChange(ctx context.Context, id uuid.UUID) bool {
	if tx, err := s.GetByID(ctx, id); err != nil {
		return false
	} else {
		switch tx.Status {
		case models.Pending:
			return true
		case models.Completed, models.Failed:
			return false
		}
	}
	return false
}
