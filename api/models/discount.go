package models

import (
	"github.com/google/uuid"
	"payment/pkg/db"
	"time"
)

type DiscountType string

const (
	Voucher DiscountType = "voucher"
	Charge  DiscountType = "charge"
)

type Discount struct {
	db.StrictBaseModel
	Code           string                 `json:"code" gorm:"not null, unique" generator:"required"`
	Description    string                 `json:"description" gorm:"not null;type:varchar(255)"`
	Amount         int64                  `json:"amount" generator:"gte=0" gorm:"not null default 0;type:integer"`
	UsageLimit     int64                  `json:"usage_limit" generator:"required,gt=0"`
	ExpirationTime time.Time              `json:"expiration_time" gorm:"not null" generator:"required"`
	Type           DiscountType           `json:"type" generator:"required" gorm:"not null;type:discount_type"`
	Transactions   []*DiscountTransaction `json:"transactions,omitempty" gorm:"foreignKey:DiscountID"`
}

func (Discount) TableName() string {
	return "discounts"
}

type DiscountResponse struct {
	Code        string       `json:"code"`
	Description string       `json:"description"`
	Total       int64        `json:"total"`
	Type        DiscountType `json:"type"`
}

type DiscountTransaction struct {
	db.StrictBaseModel
	DiscountID uuid.UUID `json:"-" gorm:"type:uuid;not null" generator:"required"`
	WalletID   uuid.UUID `json:"wallet_id" gorm:"type:uuid;not null" generator:"required"`
	PhoneNum   string    `json:"phone" generator:"required,mobile"`
}

func (DiscountTransaction) TableName() string {
	return "discount_transactions"
}

type DiscountApplyRequest struct {
	Code     string    `json:"code"`
	WalletID uuid.UUID `json:"wallet_id"`
	PhoneNum string    `json:"phone_num"`
}

type DiscountUsageRequest struct {
	Code string
}
