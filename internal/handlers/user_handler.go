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

// UserHandler provides endpoints for general user management tasks.
type UserHandler struct {
	userService service.UserService
	_           *models.User
}

// NewUserHandler creates a new instance of UserHandler with the specified service.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// FindAllUsers godoc
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
// @Failure      500     {object}  ProblemDetails
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

// FindUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve details for a single user using their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
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

// DeleteUser godoc
// @Summary      Delete user
// @Description  Permanently remove a user from the system by their unique ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  ProblemDetails
// @Failure      404  {object}  ProblemDetails
// @Failure      500  {object}  ProblemDetails
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
