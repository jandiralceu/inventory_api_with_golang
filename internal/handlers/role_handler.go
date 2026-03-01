package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

// RoleHandler organizes endpoints for creating and managing roles.
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler initializes a RoleHandler with its required dependencies.
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole godoc
// @Summary      Create a role
// @Description  Create a new role in the system
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateRoleRequest true "Role data"
// @Success      201 {object} models.Role
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      409 {object} ProblemDetails "Conflict"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal error"
// @Security     Bearer
// @Router       /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	createdRole, err := h.roleService.Create(c.Request.Context(), role)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdRole)
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Remove a role by its unique UUID
// @Tags         roles
// @Produce      json
// @Param        id path string true "Role ID (UUID)"
// @Success      204 "No content"
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal error"
// @Security     Bearer
// @Router       /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.roleService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// FindRoleByID godoc
// @Summary      Get role by ID
// @Description  Retrieve a single role by its unique UUID
// @Tags         roles
// @Produce      json
// @Param        id path string true "Role ID (UUID)"
// @Success      200 {object} models.Role
// @Failure      400 {object} ProblemDetails "Bad request"
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      404 {object} ProblemDetails "Not found"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Security     Bearer
// @Router       /roles/{id} [get]
func (h *RoleHandler) FindRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	role, err := h.roleService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, role)
}

// FindAllRoles godoc
// @Summary      Get all roles
// @Description  Retrieve a list of all roles
// @Tags         roles
// @Produce      json
// @Success      200 {array} models.Role
// @Failure      401 {object} ProblemDetails "Unauthorized"
// @Failure      403 {object} ProblemDetails "Forbidden"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal error"
// @Security     Bearer
// @Router       /roles [get]
func (h *RoleHandler) FindAllRoles(c *gin.Context) {
	roles, err := h.roleService.FindAll(c.Request.Context())
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, roles)
}
