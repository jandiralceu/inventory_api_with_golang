package models

import (
	"time"

	"github.com/google/uuid"
)

// InventoryTransaction represents a historical record of any stock movement.
// It serves as an audit log to track how, when, and why stock levels changed.
type InventoryTransaction struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Foreign Keys
	InventoryID uuid.UUID  `gorm:"type:uuid;not null;index" json:"inventory_id"`
	ProductID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"product_id"`
	WarehouseID uuid.UUID  `gorm:"type:uuid;not null;index" json:"warehouse_id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`

	// Quantities
	QuantityChange  int `gorm:"not null" json:"quantity_change"`
	QuantityBalance int `gorm:"not null" json:"quantity_balance"`

	// Classification
	TransactionType string `gorm:"type:varchar(50);not null" json:"transaction_type"`
	ReferenceID     string `gorm:"type:varchar(100)" json:"reference_id,omitempty"`
	Reason          string `gorm:"type:text" json:"reason,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"not null;default:NOW()" json:"created_at"`

	// Relationships
	Inventory *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
