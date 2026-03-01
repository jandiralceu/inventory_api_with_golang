package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

// CategoryHandler manages category-related API endpoints.
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler initializes a CategoryHandler with the provided service.
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory Godoc
// @Summary      Create a category
// @Description  Register a new product category in the system.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCategoryRequest true "Category data"
// @Success      201 {object} models.Category
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory Godoc
// @Summary      Update a category
// @Description  Modify an existing category details by ID.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id path string true "Category ID (UUID)"
// @Param        request body dto.UpdateCategoryRequest true "Category update data"
// @Success      200 {object} models.Category
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// FindAllCategories Godoc
// @Summary      List categories
// @Description  Retrieve a paginated list of categories with optional filters.
// @Tags         categories
// @Produce      json
// @Param        name query string false "Filter by name (partial match)"
// @Param        slug query string false "Filter by exact slug"
// @Param        parentId query string false "Filter by parent category ID"
// @Param        isActive query boolean false "Filter by active status"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        sort query string false "Sort field"
// @Param        order query string false "Sort order (asc/desc)"
// @Success      200 {object} dto.CategoryListResponse
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /categories [get]
func (h *CategoryHandler) FindAllCategories(c *gin.Context) {
	var req dto.GetCategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	resp, err := h.categoryService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// FindCategoryByID Godoc
// @Summary      Get category by ID
// @Description  Retrieve a single category by its unique ID.
// @Tags         categories
// @Produce      json
// @Param        id path string true "Category ID (UUID)"
// @Success      200 {object} models.Category
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /categories/{id} [get]
func (h *CategoryHandler) FindCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	category, err := h.categoryService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory Godoc
// @Summary      Delete a category
// @Description  Remove a category from the system by its ID.
// @Tags         categories
// @Produce      json
// @Param        id path string true "Category ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Security     Bearer
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
