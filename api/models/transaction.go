package models

import (
	"github.com/google/uuid"
	"payment/pkg/db"
)

type Type string
type Status string

const (
	Withdrawal Type = "withdrawal"
	Deposit    Type = "deposit"
)

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)

type Transaction struct {
	db.StrictBaseModel
	WalletID    uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	Type        Type      `json:"type" gorm:"not null;type:transaction_type"`
	Amount      int64     `json:"amount"`
	Status      Status    `json:"status" gorm:"not null;type:transaction_status"`
	Description string    `json:"description" validate:"required,description"`
	Wallet      Wallet    `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE;" json:"-"`
}

type NewTransaction struct {
	Phone       string `json:"phone"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Type        Type   `json:"type"`
}
