package wallets

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"payment/api/models"
	"payment/internal/transactions"
	"payment/pkg/db"
)

type IWallet interface {
	Create(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
	Delete(ctx context.Context, phone string) (*models.Wallet, error)
	Transaction(ctx context.Context, wallet *models.Wallet, transaction *models.Transaction) error
	GetByPhone(ctx context.Context, number string) (*models.Wallet, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
}

type WalletService struct {
	transaction transactions.ITransaction
	logger      *log.Logger
	db          *db.DB
}

func NewWallet(logger *log.Logger, db *db.DB) IWallet {
	return &WalletService{transactions.NewTransactionsService(logger, db), logger, db}
}

func (r *WalletService) Create(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {

	if err := r.db.WithContext(ctx).Save(wallet).Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}
	r.logger.WithFields(log.Fields{
		"section": "wallet",
		"mode":    "insert",
		"type":    "database",
		"wallet":  wallet.ID,
	}).Info("wallet created")
	return wallet, nil
}

func (r *WalletService) Delete(ctx context.Context, phone string) (*models.Wallet, error) {
	var wallet *models.Wallet

	wallet, err := r.GetByPhone(ctx, phone)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	if err = r.db.Unscoped().WithContext(ctx).Delete(&wallet).Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}

	r.logger.WithFields(log.Fields{
		"section": "wallet",
		"mode":    "delete",
		"type":    "database",
		"wallet":  wallet.ID,
	}).Info("wallet deleted")

	return wallet, nil
}

func (r *WalletService) Update(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error) {
	var err error
	if err = r.db.Model(&wallet).WithContext(ctx).
		Where("id = ?", wallet.ID).
		Update("amount", wallet.Amount).Error; err != nil {
		r.logger.Error(err)
		return nil, errors.New("could not update wallet balance")
	}
	return wallet, nil
}

func (r *WalletService) Transaction(ctx context.Context, wallet *models.Wallet, transaction *models.Transaction) error {
	var err error
	transaction.Status = models.Pending

	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	switch transaction.Type {
	case models.Deposit:
		wallet.Amount += transaction.Amount
	case models.Withdrawal:
		wallet.Amount -= transaction.Amount
	}

	if err = tx.Create(&transaction).Error; err != nil {
		r.logger.Error(err)
		tx.Rollback()
		return fmt.Errorf("could not create Transaction: %w", err)
	}

	// Update the wallet balance
	if _, err = r.Update(ctx, wallet); err != nil {
		if err = r.transaction.ChangeStatus(ctx, transaction.ID, models.Failed); err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update transaction status: %w", err)
		}
		return fmt.Errorf("could not update wallet balance: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		r.logger.Error(err)
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	if err = r.transaction.ChangeStatus(ctx, transaction.ID, models.Completed); err != nil {
		r.logger.Error(err)
		return fmt.Errorf("could not update Transaction status: %w", err)
	}

	r.logger.WithFields(log.Fields{
		"section": "transaction",
		"mode":    "insert",
		"type":    "database",
		"transaction": map[string]interface{}{
			"id":          transaction.ID.String(),
			"type":        transaction.Type,
			"description": transaction.Description,
		},
	}).Info("transaction successfully created")

	return nil
}

func (r *WalletService) GetByPhone(ctx context.Context, number string) (*models.Wallet, error) {
	var wallet *models.Wallet
	var err error

	if err = r.db.WithContext(ctx).First(&wallet, "phone = ?", number).Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}
	if wallet.Transactions, err = r.transaction.List(ctx, wallet.ID); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return wallet, nil
}

func (r *WalletService) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet *models.Wallet
	var err error

	if err = r.db.WithContext(ctx).First(&wallet, "id = ?", id).Error; err != nil {
		r.logger.Error(err)
		return nil, err
	}
	if wallet.Transactions, err = r.transaction.List(ctx, wallet.ID); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return wallet, nil
}
