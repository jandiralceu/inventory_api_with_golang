package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct Godoc
// @Summary      Create a product
// @Description  Register a new product in the system.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductRequest true "Product data"
// @Success      201 {object} models.Product
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	product, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct Godoc
// @Summary      Update a product
// @Description  Modify an existing product details by ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path string true "Product ID (UUID)"
// @Param        request body dto.UpdateProductRequest true "Product update data"
// @Success      200 {object} models.Product
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// FindAllProducts Godoc
// @Summary      List products
// @Description  Retrieve a paginated list of products with optional filters.
// @Tags         products
// @Produce      json
// @Param        name query string false "Filter by name (partial match)"
// @Param        sku query string false "Filter by exact SKU"
// @Param        categoryId query string false "Filter by Category ID"
// @Param        supplierId query string false "Filter by Supplier ID"
// @Param        isActive query boolean false "Filter by active status"
// @Param        minPrice query number false "Filter by minimum price"
// @Param        maxPrice query number false "Filter by maximum price"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.ProductListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /products [get]
func (h *ProductHandler) FindAllProducts(c *gin.Context) {
	var req dto.GetProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.productService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FindProductByID Godoc
// @Summary      Get product by ID
// @Description  Retrieve a single product by its unique ID.
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID (UUID)"
// @Success      200 {object} models.Product
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /products/{id} [get]
func (h *ProductHandler) FindProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	product, err := h.productService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct Godoc
// @Summary      Delete a product
// @Description  Remove a product from the system by its ID.
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
