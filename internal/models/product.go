package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProductImages represents an array of product image URLs
type ProductImages []string

// Value makes ProductImages implement driver.Valuer for JSONB storage
func (pi ProductImages) Value() (driver.Value, error) {
	if pi == nil {
		return nil, nil
	}
	return json.Marshal(pi)
}

// Scan makes ProductImages implement sql.Scanner for JSONB retrieval
func (pi *ProductImages) Scan(value any) error {
	if value == nil {
		*pi = ProductImages{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &pi)
}

// ProductMetadata represents flexible JSON attributes
type ProductMetadata map[string]any

// Value makes ProductMetadata implement driver.Valuer for JSONB storage
func (pm ProductMetadata) Value() (driver.Value, error) {
	if pm == nil {
		return nil, nil
	}
	return json.Marshal(pm)
}

// Scan makes ProductMetadata implement sql.Scanner for JSONB retrieval
func (pm *ProductMetadata) Scan(value any) error {
	if value == nil {
		*pm = ProductMetadata{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &pm)
}

// Product represents the product entity in the database
type Product struct {
	ID              uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	SKU             string          `gorm:"type:varchar(50);not null;uniqueIndex" json:"sku"`
	Slug            string          `gorm:"type:varchar(200);not null;uniqueIndex" json:"slug"`
	Name            string          `gorm:"type:varchar(200);not null;index" json:"name"`
	Description     string          `gorm:"type:text" json:"description,omitempty"`
	Price           float64         `gorm:"type:decimal(10,2);not null;index" json:"price"`
	CostPrice       *float64        `gorm:"type:decimal(10,2)" json:"costPrice,omitempty"`
	CategoryID      *uuid.UUID      `gorm:"type:uuid;index" json:"categoryId,omitempty"`
	SupplierID      *uuid.UUID      `gorm:"type:uuid;index" json:"supplierId,omitempty"`
	IsActive        bool            `gorm:"not null;default:true;index" json:"isActive"`
	ReorderLevel    int             `gorm:"not null;default:0" json:"reorderLevel"`
	ReorderQuantity int             `gorm:"not null;default:0" json:"reorderQuantity"`
	WeightKg        *float64        `gorm:"type:decimal(10,3)" json:"weightKg,omitempty"`
	Images          ProductImages   `gorm:"type:jsonb;index:,type:gin" json:"images,omitempty"`
	Metadata        ProductMetadata `gorm:"type:jsonb;index:,type:gin" json:"metadata,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Supplier *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}
