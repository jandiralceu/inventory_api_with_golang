package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// =====================
// PaginationRequest Tests
// =====================

func TestGetPageDefault(t *testing.T) {
	req := &PaginationRequest{}
	assert.Equal(t, 1, req.GetPage())
}

func TestGetPageZero(t *testing.T) {
	req := &PaginationRequest{Page: 0}
	assert.Equal(t, 1, req.GetPage())
}

func TestGetPageNegative(t *testing.T) {
	req := &PaginationRequest{Page: -5}
	assert.Equal(t, 1, req.GetPage())
}

func TestGetPageCustom(t *testing.T) {
	req := &PaginationRequest{Page: 3}
	assert.Equal(t, 3, req.GetPage())
}

func TestGetLimitDefault(t *testing.T) {
	req := &PaginationRequest{}
	assert.Equal(t, 10, req.GetLimit())
}

func TestGetLimitZero(t *testing.T) {
	req := &PaginationRequest{Limit: 0}
	assert.Equal(t, 10, req.GetLimit())
}

func TestGetLimitNegative(t *testing.T) {
	req := &PaginationRequest{Limit: -1}
	assert.Equal(t, 10, req.GetLimit())
}

func TestGetLimitCustom(t *testing.T) {
	req := &PaginationRequest{Limit: 50}
	assert.Equal(t, 50, req.GetLimit())
}

func TestGetSortDefault(t *testing.T) {
	req := &PaginationRequest{}
	assert.Equal(t, "created_at", req.GetSort("name", "email"))
}

func TestGetSortAllowedField(t *testing.T) {
	req := &PaginationRequest{Sort: "title"}
	assert.Equal(t, "title", req.GetSort("title", "name"))
}

func TestGetSortDisallowedField(t *testing.T) {
	req := &PaginationRequest{Sort: "password_hash"}
	assert.Equal(t, "created_at", req.GetSort("name", "email"))
}

func TestGetSortNoAllowedFields(t *testing.T) {
	req := &PaginationRequest{Sort: "name"}
	assert.Equal(t, "created_at", req.GetSort())
}

func TestGetOrderDefault(t *testing.T) {
	req := &PaginationRequest{}
	assert.Equal(t, "desc", req.GetOrder())
}

func TestGetOrderCustom(t *testing.T) {
	req := &PaginationRequest{Order: "asc"}
	assert.Equal(t, "asc", req.GetOrder())
}

// =====================
// NewPaginatedResponse Tests
// =====================

func TestNewPaginatedResponseSinglePage(t *testing.T) {
	data := []string{"a", "b", "c"}
	resp := NewPaginatedResponse(data, 3, 1, 10)

	assert.Equal(t, 3, len(resp.Data))
	assert.Equal(t, int64(3), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 1, resp.TotalPages)
}

func TestNewPaginatedResponseMultiplePages(t *testing.T) {
	data := []string{"a", "b"}
	resp := NewPaginatedResponse(data, 25, 2, 10)

	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, int64(25), resp.Total)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 3, resp.TotalPages)
}

func TestNewPaginatedResponseExactDivision(t *testing.T) {
	data := []string{"a"}
	resp := NewPaginatedResponse(data, 20, 1, 10)

	assert.Equal(t, 2, resp.TotalPages)
}

func TestNewPaginatedResponseEmptyData(t *testing.T) {
	data := []string{}
	resp := NewPaginatedResponse(data, 0, 1, 10)

	assert.Equal(t, 0, len(resp.Data))
	assert.Equal(t, int64(0), resp.Total)
	assert.Equal(t, 0, resp.TotalPages)
}

func TestNewPaginatedResponseLimitOne(t *testing.T) {
	data := []string{"a"}
	resp := NewPaginatedResponse(data, 5, 1, 1)

	assert.Equal(t, 5, resp.TotalPages)
}
