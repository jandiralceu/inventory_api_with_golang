package dto

import (
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// SupplierAddress matches the address structure for API requests/responses.
type SupplierAddress struct {
	Street     string `json:"street" binding:"required,min=3"`
	Number     string `json:"number" binding:"required"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required,len=2"`
	Country    string `json:"country" binding:"required"`
	ZipCode    string `json:"zipCode" binding:"required"`
}

// MapToModel converts the DTO address to the domain model address.
func (a SupplierAddress) MapToModel() models.Address {
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

// CreateSupplierRequest defines the payload for registering a new supplier.
type CreateSupplierRequest struct {
	Name          string          `json:"name" binding:"required,min=3,max=100"`
	Description   string          `json:"description" binding:"omitempty"`
	TaxID         string          `json:"taxId" binding:"required,min=5,max=50"`
	Email         string          `json:"email" binding:"omitempty,email"`
	Phone         string          `json:"phone" binding:"omitempty"`
	Address       SupplierAddress `json:"address" binding:"required"`
	ContactPerson string          `json:"contactPerson" binding:"omitempty"`
}

// UpdateSupplierRequest defines the payload for modifying supplier data.
type UpdateSupplierRequest struct {
	Name          string          `json:"name" binding:"required,min=3,max=100"`
	Description   string          `json:"description" binding:"omitempty"`
	TaxID         string          `json:"taxId" binding:"required,min=5,max=50"`
	Email         string          `json:"email" binding:"omitempty,email"`
	Phone         string          `json:"phone" binding:"omitempty"`
	Address       SupplierAddress `json:"address" binding:"required"`
	ContactPerson string          `json:"contactPerson" binding:"omitempty"`
	IsActive      *bool           `json:"isActive" binding:"required"`
}

// GetSupplierListRequest defines filters and pagination parameters for suppliers.
type GetSupplierListRequest struct {
	PaginationRequest
	Name     string `form:"name" binding:"omitempty"`
	TaxID    string `form:"taxId" binding:"omitempty"`
	Email    string `form:"email" binding:"omitempty"`
	IsActive *bool  `form:"isActive" binding:"omitempty"`
}

// SupplierListResponse matches the paginated structure for supplier collections.
type SupplierListResponse PaginatedResponse[models.Supplier]
