package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// TransactionListRequest defines query parameters for listing inventory transactions (audit log).
type TransactionListRequest struct {
	PaginationRequest
	InventoryID     *uuid.UUID `form:"inventoryId"`
	ProductID       *uuid.UUID `form:"productId"`
	WarehouseID     *uuid.UUID `form:"warehouseId"`
	UserID          *uuid.UUID `form:"userId"`
	TransactionType string     `form:"transactionType"`
	StartDate       *time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate         *time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`
}

// TransactionResponse represents a single stock movement record in the audit log.
type TransactionResponse struct {
	ID              uuid.UUID         `json:"id"`
	InventoryID     uuid.UUID         `json:"inventoryId"`
	ProductID       uuid.UUID         `json:"productId"`
	WarehouseID     uuid.UUID         `json:"warehouseId"`
	UserID          *uuid.UUID        `json:"userId,omitempty"`
	QuantityChange  int               `json:"quantityChange"`
	QuantityBalance int               `json:"quantityBalance"`
	TransactionType string            `json:"transactionType"`
	ReferenceID     string            `json:"referenceId,omitempty"`
	Reason          string            `json:"reason,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
	Product         *models.Product   `json:"product,omitempty"`
	Warehouse       *models.Warehouse `json:"warehouse,omitempty"`
	User            *models.User      `json:"user,omitempty"`
}

// TransactionListResponse represents a paginated list of audit records.
type TransactionListResponse struct {
	PaginatedResponse[TransactionResponse]
}
