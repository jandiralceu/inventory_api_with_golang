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

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

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

func (h *RoleHandler) FindAllRoles(c *gin.Context) {
	roles, err := h.roleService.FindAll(c.Request.Context())
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, roles)
}
