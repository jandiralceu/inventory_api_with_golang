package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

// WarehouseHandler manages warehouse-related API endpoints.
type WarehouseHandler struct {
	warehouseService service.WarehouseService
}

// NewWarehouseHandler initializes a WarehouseHandler with the provided service.
func NewWarehouseHandler(warehouseService service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

// CreateWarehouse Godoc
// @Summary      Create a warehouse
// @Description  Register a new warehouse in the system.
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateWarehouseRequest true "Warehouse data"
// @Success      201 {object} models.Warehouse
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /warehouses [post]
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req dto.CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	warehouse, err := h.warehouseService.Create(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, warehouse)
}

// UpdateWarehouse Godoc
// @Summary      Update a warehouse
// @Description  Modify an existing warehouse details by ID.
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id path string true "Warehouse ID (UUID)"
// @Param        request body dto.UpdateWarehouseRequest true "Warehouse update data"
// @Success      200 {object} models.Warehouse
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /warehouses/{id} [put]
func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	warehouse, err := h.warehouseService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// FindAllWarehouses Godoc
// @Summary      List warehouses
// @Description  Retrieve a paginated list of warehouses with optional filters.
// @Tags         warehouses
// @Produce      json
// @Param        name query string false "Filter by name (partial match)"
// @Param        code query string false "Filter by exact code"
// @Param        isActive query boolean false "Filter by active status"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.WarehouseListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /warehouses [get]
func (h *WarehouseHandler) FindAllWarehouses(c *gin.Context) {
	var req dto.GetWarehouseListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.warehouseService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FindWarehouseByID Godoc
// @Summary      Get warehouse by ID
// @Description  Retrieve a single warehouse by its unique ID.
// @Tags         warehouses
// @Produce      json
// @Param        id path string true "Warehouse ID (UUID)"
// @Success      200 {object} models.Warehouse
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /warehouses/{id} [get]
func (h *WarehouseHandler) FindWarehouseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	warehouse, err := h.warehouseService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouse Godoc
// @Summary      Delete a warehouse
// @Description  Remove a warehouse from the system by its ID.
// @Tags         warehouses
// @Produce      json
// @Param        id path string true "Warehouse ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /warehouses/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.warehouseService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
