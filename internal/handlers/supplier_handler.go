package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

// SupplierHandler manages supplier-related API endpoints.
type SupplierHandler struct {
	supplierService service.SupplierService
}

// NewSupplierHandler initializes a SupplierHandler with the provided service.
func NewSupplierHandler(supplierService service.SupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// CreateSupplier Godoc
// @Summary      Create a supplier
// @Description  Register a new supplier in the system.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSupplierRequest true "Supplier data"
// @Success      201 {object} models.Supplier
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /suppliers [post]
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req dto.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	supplier, err := h.supplierService.Create(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// UpdateSupplier Godoc
// @Summary      Update a supplier
// @Description  Modify an existing supplier details by ID.
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        id path string true "Supplier ID (UUID)"
// @Param        request body dto.UpdateSupplierRequest true "Supplier update data"
// @Success      200 {object} models.Supplier
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	supplier, err := h.supplierService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// FindAllSuppliers Godoc
// @Summary      List suppliers
// @Description  Retrieve a paginated list of suppliers with optional filters.
// @Tags         suppliers
// @Produce      json
// @Param        name query string false "Filter by name (partial match)"
// @Param        taxId query string false "Filter by exact Tax ID"
// @Param        email query string false "Filter by email"
// @Param        isActive query boolean false "Filter by active status"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.SupplierListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /suppliers [get]
func (h *SupplierHandler) FindAllSuppliers(c *gin.Context) {
	var req dto.GetSupplierListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.supplierService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FindSupplierByID Godoc
// @Summary      Get supplier by ID
// @Description  Retrieve a single supplier by its unique ID.
// @Tags         suppliers
// @Produce      json
// @Param        id path string true "Supplier ID (UUID)"
// @Success      200 {object} models.Supplier
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /suppliers/{id} [get]
func (h *SupplierHandler) FindSupplierByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	supplier, err := h.supplierService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier Godoc
// @Summary      Delete a supplier
// @Description  Remove a supplier from the system by its ID.
// @Tags         suppliers
// @Produce      json
// @Param        id path string true "Supplier ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.supplierService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
