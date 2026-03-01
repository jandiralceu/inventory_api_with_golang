package dto

import (
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// CreateCategoryRequest defines the payload for creating a new product category.
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required,min=3,max=100"`
	Description string     `json:"description" binding:"omitempty"`
	ParentID    *uuid.UUID `json:"parentId" binding:"omitempty"`
}

// UpdateCategoryRequest defines the payload for modifying an existing category.
type UpdateCategoryRequest struct {
	Name        string     `json:"name" binding:"required,min=3,max=100"`
	Description string     `json:"description" binding:"omitempty"`
	ParentID    *uuid.UUID `json:"parentId" binding:"omitempty"`
	IsActive    *bool      `json:"isActive" binding:"required"`
}

// GetCategoryListRequest defines filters and pagination for categories retrieval.
type GetCategoryListRequest struct {
	PaginationRequest
	Name     string     `form:"name" binding:"omitempty"`
	Slug     string     `form:"slug" binding:"omitempty"`
	ParentID *uuid.UUID `form:"parentId" binding:"omitempty"`
	IsActive *bool      `form:"isActive" binding:"omitempty"`
}

// CategoryListResponse matches the paginated structure for category collections.
type CategoryListResponse PaginatedResponse[models.Category]
