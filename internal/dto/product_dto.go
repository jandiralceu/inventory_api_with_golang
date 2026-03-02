package dto

import (
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// CreateProductRequest defines the payload for creating a product
type CreateProductRequest struct {
	SKU             string                 `json:"sku" binding:"required,min=2,max=50"`
	Name            string                 `json:"name" binding:"required,min=2,max=200"`
	Description     string                 `json:"description,omitempty"`
	Price           float64                `json:"price" binding:"required,min=0"`
	CostPrice       *float64               `json:"costPrice,omitempty" binding:"omitempty,min=0"`
	CategoryID      *uuid.UUID             `json:"categoryId,omitempty"`
	SupplierID      *uuid.UUID             `json:"supplierId,omitempty"`
	ReorderLevel    int                    `json:"reorderLevel,omitempty" binding:"min=0"`
	ReorderQuantity int                    `json:"reorderQuantity,omitempty" binding:"min=0"`
	WeightKg        *float64               `json:"weightKg,omitempty" binding:"omitempty,min=0"`
	Images          []string               `json:"images,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateProductRequest defines the payload for updating a product
type UpdateProductRequest struct {
	SKU             string                 `json:"sku,omitempty" binding:"omitempty,min=2,max=50"`
	Name            string                 `json:"name,omitempty" binding:"omitempty,min=2,max=200"`
	Description     string                 `json:"description,omitempty"`
	Price           *float64               `json:"price,omitempty" binding:"omitempty,min=0"`
	CostPrice       *float64               `json:"costPrice,omitempty" binding:"omitempty,min=0"`
	CategoryID      *uuid.UUID             `json:"categoryId,omitempty"`
	SupplierID      *uuid.UUID             `json:"supplierId,omitempty"`
	ReorderLevel    *int                   `json:"reorderLevel,omitempty" binding:"omitempty,min=0"`
	ReorderQuantity *int                   `json:"reorderQuantity,omitempty" binding:"omitempty,min=0"`
	WeightKg        *float64               `json:"weightKg,omitempty" binding:"omitempty,min=0"`
	Images          []string               `json:"images,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	IsActive        *bool                  `json:"isActive,omitempty"`
}

// GetProductListRequest defines query parameters for listing products
type GetProductListRequest struct {
	PaginationRequest
	Name       string     `form:"name"`
	SKU        string     `form:"sku"`
	CategoryID *uuid.UUID `form:"categoryId"`
	SupplierID *uuid.UUID `form:"supplierId"`
	IsActive   *bool      `form:"isActive"`
	MinPrice   *float64   `form:"minPrice"`
	MaxPrice   *float64   `form:"maxPrice"`
}

// ProductListResponse represents the paginated response for products
type ProductListResponse struct {
	PaginatedResponse[models.Product]
}
