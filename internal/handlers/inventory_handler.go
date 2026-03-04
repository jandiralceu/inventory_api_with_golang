package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory Godoc
// @Summary      Create inventory record
// @Description  Initialize a new inventory record for a product in a warehouse.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateInventoryRequest true "Inventory data"
// @Success      201 {object} models.Inventory
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	inventory, err := h.inventoryService.Create(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

// UpdateInventory Godoc
// @Summary      Update inventory record
// @Description  Update details like location code or inventory levels.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Param        request body dto.UpdateInventoryRequest true "Inventory update data"
// @Success      200 {object} models.Inventory
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id} [put]
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	inventory, err := h.inventoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// FindAllInventory Godoc
// @Summary      List inventory records
// @Description  Retrieve a paginated list of inventory records with optional filters.
// @Tags         inventories
// @Produce      json
// @Param        productId query string false "Filter by Product ID"
// @Param        warehouseId query string false "Filter by Warehouse ID"
// @Param        lowStock query boolean false "Filter items with low stock"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.InventoryListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories [get]
func (h *InventoryHandler) FindAllInventory(c *gin.Context) {
	var req dto.GetInventoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.inventoryService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FindInventoryByID Godoc
// @Summary      Get inventory by ID
// @Description  Retrieve a single inventory record by its unique ID.
// @Tags         inventories
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Success      200 {object} models.Inventory
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id} [get]
func (h *InventoryHandler) FindInventoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	inventory, err := h.inventoryService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// AddStock Godoc
// @Summary      Inbound stock
// @Description  Add physical stock to a specific inventory record.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Param        request body dto.StockOperationRequest true "Stock operation data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id}/add [post]
func (h *InventoryHandler) AddStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	if err := h.inventoryService.AddStock(c.Request.Context(), id, req.Quantity); err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock added successfully"})
}

// RemoveStock Godoc
// @Summary      Outbound stock
// @Description  Remove physical stock from a specific inventory record.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Param        request body dto.StockOperationRequest true "Stock operation data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id}/remove [post]
func (h *InventoryHandler) RemoveStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	if err := h.inventoryService.RemoveStock(c.Request.Context(), id, req.Quantity); err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock removed successfully"})
}

// ReserveStock Godoc
// @Summary      Reserve stock
// @Description  Reserve a quantity of stock for a pending order.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Param        request body dto.StockOperationRequest true "Stock operation data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id}/reserve [post]
func (h *InventoryHandler) ReserveStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	if err := h.inventoryService.ReserveStock(c.Request.Context(), id, req.Quantity); err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock reserved successfully"})
}

// ReleaseStock Godoc
// @Summary      Release reserved stock
// @Description  Release a previously reserved stock quantity.
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Param        request body dto.StockOperationRequest true "Stock operation data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id}/release [post]
func (h *InventoryHandler) ReleaseStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.StockOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	if err := h.inventoryService.ReleaseStock(c.Request.Context(), id, req.Quantity); err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock released successfully"})
}

// DeleteInventory Godoc
// @Summary      Delete inventory record
// @Description  Remove an inventory record from the system.
// @Tags         inventories
// @Produce      json
// @Param        id path string true "Inventory ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/{id} [delete]
func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.inventoryService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTransactionHistory Godoc
// @Summary      Stock movement history
// @Description  Retrieve a global audit log of all stock movements with advanced filters.
// @Tags         inventories
// @Produce      json
// @Param        inventoryId query string false "Filter by Inventory ID"
// @Param        productId query string false "Filter by Product ID"
// @Param        warehouseId query string false "Filter by Warehouse ID"
// @Param        userId query string false "Filter by User ID"
// @Param        transactionType query string false "Filter by Type"
// @Param        startDate query string false "Start date (ISO8601)"
// @Param        endDate query string false "End date (ISO8601)"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.TransactionListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /inventories/transactions [get]
func (h *InventoryHandler) GetTransactionHistory(c *gin.Context) {
	var req dto.TransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.inventoryService.GetTransactionHistory(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
