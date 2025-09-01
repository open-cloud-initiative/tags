package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag ...
type Tag struct {
	// ID ...
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	// Name is the tag name.
	Name string `json:"name" gorm:"uniqueIndex:idx_account_name_value"`
	// Value is the tag value.
	Value string `json:"value" gorm:"uniqueIndex:idx_account_name_value"`
	// OwnedBy ...
	OwnedBy uuid.UUID `json:"ownedBy" gorm:"uniqueIndex:idx_account_name_value;type:uuid;not null"`
	// CreatedAt ...
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt ...
	UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt ...
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
