package models

import "payment/pkg/db"

type Wallet struct {
	db.BaseModel
	Phone        string         `gorm:"unique;type:varchar(20)" json:"phone,omitempty"`
	Amount       int64          `gorm:"type:int;default:0" json:"amount,omitempty"`
	Transactions []*Transaction `gorm:"constraint:OnDelete:CASCADE;" json:"transactions,omitempty"`
}
