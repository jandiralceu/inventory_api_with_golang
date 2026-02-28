package dto

import "math"

// PaginationRequest encapsulates common parameters for paginated list requests.
type PaginationRequest struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Order string `form:"order" binding:"omitempty,oneof=asc desc"`
	Sort  string `form:"sort" binding:"omitempty"`
}

func (p *PaginationRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit <= 0 {
		return 10 // default limit
	}
	return p.Limit
}

func (p *PaginationRequest) GetSort(allowedFields ...string) string {
	if p.Sort != "" {
		for _, field := range allowedFields {
			if p.Sort == field {
				return p.Sort
			}
		}
	}
	return "created_at"
}

func (p *PaginationRequest) GetOrder() string {
	if p.Order == "" {
		return "desc"
	}
	return p.Order
}

// PaginatedResponse is a generic wrapper for API responses that return a list of items.
type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}

// NewPaginatedResponse computes the total pages and return a [PaginatedResponse].
func NewPaginatedResponse[T any](data []T, total int64, page, limit int) PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return PaginatedResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
