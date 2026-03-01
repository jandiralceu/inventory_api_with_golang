package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

// UserHandler provides endpoints for general user management tasks.
type UserHandler struct {
	userService service.UserService
	_           *models.User
}

// NewUserHandler creates a new instance of UserHandler with the specified service for user operations.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// FindAllUsers Godoc
// @Summary      List all users
// @Description  Get a paginated list of users. Supports filtering by name, email, and role.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page    query  int     false  "Page number"
// @Param        limit   query  int     false  "Number of items per page"
// @Param        name    query  string  false  "Filter by name"
// @Param        email   query  string  false  "Filter by email"
// @Param        roleId  query  string  false  "Filter by role UUID"
// @Success      200     {object}  dto.UserListResponse
// @Failure      400     {object}  ProblemDetails
// @Failure      401     {object}  ProblemDetails
// @Failure      403     {object}  ProblemDetails
// @Failure      429     {object}  ProblemDetails
// @Failure      500     {object}  ProblemDetails
// @Security     Bearer
// @Router       /users [get]
func (h *UserHandler) FindAllUsers(c *gin.Context) {
	var req dto.GetUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	users, err := h.userService.FindAll(c.Request.Context(), req)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// FindUserByID Godoc
// @Summary      Get user by ID
// @Description  Retrieve details for a single user using their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/{id} [get]
func (h *UserHandler) FindUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	user, err := h.userService.FindByID(c.Request.Context(), id)
	if err != nil {
		RespondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser Godoc
// @Summary      Delete user
// @Description  Permanently remove a user from the system by their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(c, apperrors.ErrInvalidID)
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangePassword Godoc
// @Summary      Change own password
// @Description  Updates the authenticated user's password. Requires the old password for verification.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.ChangePasswordRequest true "Password change data"
// @Success      204 "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/change-password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	userID := middleware.GetUserID(c)
	if userID == uuid.Nil {
		RespondWithError(c, apperrors.ErrUnauthorized)
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID, req); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangeRole Godoc
// @Summary      Change user role
// @Description  Updates a user's assigned role. Typically an administrative task.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.ChangeRoleRequest true "User and Role UUIDs"
// @Success      204 "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      401  {object}  ProblemDetails
// @Failure      403  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
// @Failure      429  {object}  ProblemDetails
// @Security     Bearer
// @Router       /users/change-role [patch]
func (h *UserHandler) ChangeRole(c *gin.Context) {
	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	if err := h.userService.ChangeRole(c.Request.Context(), req.UserID, req); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
