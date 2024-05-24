package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// BaseModel provides common fields for database models,
// including a unique ID, timestamps for creation and updates, and a soft delete field.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// StrictBaseModel are used for the structs that doesn't require updating the timestamp data or soft deletion.
type StrictBaseModel struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
}
