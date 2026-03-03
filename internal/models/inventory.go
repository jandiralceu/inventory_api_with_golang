package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// InventoryMetadata represents flexible JSON attributes for inventory items
type InventoryMetadata map[string]any

// Value makes InventoryMetadata implement driver.Valuer for JSONB storage
func (im InventoryMetadata) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	return json.Marshal(im)
}

// Scan makes InventoryMetadata implement sql.Scanner for JSONB retrieval
func (im *InventoryMetadata) Scan(value any) error {
	if value == nil {
		*im = InventoryMetadata{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &im)
}

// Inventory represents the inventory entity in the database
type Inventory struct {
	ID               uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID        uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:uk_inventory_product_warehouse" json:"productId"`
	WarehouseID      uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:uk_inventory_product_warehouse" json:"warehouseId"`
	Quantity         int               `gorm:"not null;default:0;check:quantity >= 0" json:"quantity"`
	ReservedQuantity int               `gorm:"not null;default:0;check:reserved_quantity >= 0" json:"reservedQuantity"`
	LocationCode     string            `gorm:"type:varchar(50)" json:"locationCode,omitempty"`
	MinQuantity      int               `gorm:"default:0;check:min_quantity >= 0" json:"minQuantity"`
	MaxQuantity      *int              `gorm:"type:integer" json:"maxQuantity,omitempty"`
	Version          int               `gorm:"not null;default:1" json:"version"`
	LastCountedAt    *time.Time        `json:"lastCountedAt,omitempty"`
	Metadata         InventoryMetadata `gorm:"type:jsonb;index:,type:gin" json:"metadata,omitempty"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
}

// TableName overrides the table name used by Inventory to `inventory`
func (Inventory) TableName() string {
	return "inventory"
}

// AvailableQuantity returns the quantity available for sale/new reservations
func (i *Inventory) AvailableQuantity() int {
	return i.Quantity - i.ReservedQuantity
}
