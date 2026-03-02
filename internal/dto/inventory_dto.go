package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// CreateInventoryRequest defines the payload for creating an inventory record
type CreateInventoryRequest struct {
	ProductID    uuid.UUID      `json:"productId" binding:"required"`
	WarehouseID  uuid.UUID      `json:"warehouseId" binding:"required"`
	Quantity     int            `json:"quantity" binding:"min=0"`
	LocationCode string         `json:"locationCode,omitempty"`
	MinQuantity  int            `json:"minQuantity,omitempty" binding:"min=0"`
	MaxQuantity  *int           `json:"maxQuantity,omitempty" binding:"omitempty,gtfield=MinQuantity"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// UpdateInventoryRequest defines the payload for updating an inventory record
type UpdateInventoryRequest struct {
	Quantity         *int           `json:"quantity,omitempty" binding:"omitempty,min=0"`
	ReservedQuantity *int           `json:"reservedQuantity,omitempty" binding:"omitempty,min=0"`
	LocationCode     string         `json:"locationCode,omitempty"`
	MinQuantity      *int           `json:"minQuantity,omitempty" binding:"omitempty,min=0"`
	MaxQuantity      *int           `json:"maxQuantity,omitempty" binding:"omitempty,min=0"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

// StockOperationRequest defines a request to add or remove stock
type StockOperationRequest struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Reason   string `json:"reason,omitempty"`
}

// InventoryResponse represents the inventory data returned to the client
type InventoryResponse struct {
	ID                uuid.UUID                `json:"id"`
	ProductID         uuid.UUID                `json:"productId"`
	WarehouseID       uuid.UUID                `json:"warehouseId"`
	Quantity          int                      `json:"quantity"`
	ReservedQuantity  int                      `json:"reservedQuantity"`
	AvailableQuantity int                      `json:"availableQuantity"`
	LocationCode      string                   `json:"locationCode,omitempty"`
	MinQuantity       int                      `json:"minQuantity"`
	MaxQuantity       *int                     `json:"maxQuantity,omitempty"`
	Version           int                      `json:"version"`
	LastCountedAt     *time.Time               `json:"lastCountedAt,omitempty"`
	Metadata          models.InventoryMetadata `json:"metadata,omitempty"`
	CreatedAt         time.Time                `json:"createdAt"`
	UpdatedAt         time.Time                `json:"updatedAt"`
	Product           *models.Product          `json:"product,omitempty"`
	Warehouse         *models.Warehouse        `json:"warehouse,omitempty"`
}

// GetInventoryListRequest defines query parameters for listing inventory
type GetInventoryListRequest struct {
	PaginationRequest
	ProductID   *uuid.UUID `form:"productId"`
	WarehouseID *uuid.UUID `form:"warehouseId"`
	LowStock    *bool      `form:"lowStock"` // Filter items where quantity <= min_quantity
}

// InventoryListResponse represents the paginated response for inventory
type InventoryListResponse struct {
	PaginatedResponse[InventoryResponse]
}
