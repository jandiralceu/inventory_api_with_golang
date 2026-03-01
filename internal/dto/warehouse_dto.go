package dto

import (
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// WarehouseAddress matches the address structure for API requests/responses.
type WarehouseAddress struct {
	Street     string `json:"street" binding:"required,min=3"`
	Number     string `json:"number" binding:"required"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required,len=2"`
	Country    string `json:"country" binding:"required"`
	ZipCode    string `json:"zipCode" binding:"required"`
}

// MapToModel converts the DTO address to the domain model address.
func (a WarehouseAddress) MapToModel() models.Address {
	return models.Address{
		Street:     a.Street,
		Number:     a.Number,
		Complement: a.Complement,
		City:       a.City,
		State:      a.State,
		Country:    a.Country,
		ZipCode:    a.ZipCode,
	}
}

// CreateWarehouseRequest defines the payload for registering a new warehouse.
type CreateWarehouseRequest struct {
	Name        string           `json:"name" binding:"required,min=3,max=100"`
	Code        string           `json:"code" binding:"required,min=2,max=50"`
	Description string           `json:"description" binding:"omitempty"`
	Address     WarehouseAddress `json:"address" binding:"required"`
	ManagerName string           `json:"managerName" binding:"omitempty"`
	Phone       string           `json:"phone" binding:"omitempty"`
	Email       string           `json:"email" binding:"omitempty,email"`
	Notes       string           `json:"notes" binding:"omitempty"`
}

// UpdateWarehouseRequest defines the payload for modifying warehouse data.
type UpdateWarehouseRequest struct {
	Name        string           `json:"name" binding:"required,min=3,max=100"`
	Code        string           `json:"code" binding:"required,min=2,max=50"`
	Description string           `json:"description" binding:"omitempty"`
	Address     WarehouseAddress `json:"address" binding:"required"`
	ManagerName string           `json:"managerName" binding:"omitempty"`
	Phone       string           `json:"phone" binding:"omitempty"`
	Email       string           `json:"email" binding:"omitempty,email"`
	Notes       string           `json:"notes" binding:"omitempty"`
	IsActive    *bool            `json:"isActive" binding:"required"`
}

// GetWarehouseListRequest defines filters and pagination parameters for warehouses.
type GetWarehouseListRequest struct {
	PaginationRequest
	Name     string `form:"name" binding:"omitempty"`
	Code     string `form:"code" binding:"omitempty"`
	IsActive *bool  `form:"isActive" binding:"omitempty"`
}

// WarehouseListResponse matches the paginated structure for warehouse collections.
type WarehouseListResponse PaginatedResponse[models.Warehouse]
