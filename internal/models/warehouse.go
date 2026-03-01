package models

import (
	"time"

	"github.com/google/uuid"
)

// Warehouse represents a physical location where products are stored.
type Warehouse struct {
	// ID is the unique identifier for the warehouse.
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Name is the display name of the warehouse.
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	// Slug is the URL-friendly name used for queries.
	Slug string `gorm:"type:varchar(100);not null;unique" json:"slug"`
	// Code is the internal reference code for the warehouse (e.g., 'WH01', 'SP01').
	Code string `gorm:"type:varchar(50);not null;unique" json:"code"`
	// Description provides additional context about the warehouse.
	Description string `gorm:"type:text" json:"description"`
	// Address is a flexible JSONB field containing location details.
	// Reuses the Address struct defined in supplier.go
	Address Address `gorm:"type:jsonb" json:"address"`
	// IsActive determines if the warehouse is currently enabled for operations.
	IsActive bool `gorm:"type:boolean;not null;default:true" json:"is_active"`
	// ManagerName is the name of the person responsible for the warehouse.
	ManagerName string `gorm:"type:varchar(100)" json:"manager_name"`
	// Phone is the contact phone number for the warehouse.
	Phone string `gorm:"type:varchar(20)" json:"phone"`
	// Email is the contact email for the warehouse.
	Email string `gorm:"type:varchar(255)" json:"email"`
	// Notes contains any extra information about the warehouse.
	Notes string `gorm:"type:text" json:"notes"`
	// CreatedAt is the timestamp when the warehouse was registered.
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	// UpdatedAt is the timestamp for the last update to the warehouse record.
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}
