package dto

import (
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
)

// ChangePasswordRequest defines the data required to rotate a user's password.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// ChangeRoleRequest defines the payload for updating a user's role assignment.
type ChangeRoleRequest struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
	RoleID uuid.UUID `json:"roleId" binding:"required"`
}

// CreateUserRequest defines the validation rules for creating a new user.
type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required,min=3,max=100"`
	Email    string    `json:"email" binding:"required,email,max=255"`
	Password string    `json:"password" binding:"required,min=8"`
	RoleID   uuid.UUID `json:"roleId" binding:"required"`
}

// GetUserListRequest defines filters and pagination for users retrieval.
type GetUserListRequest struct {
	PaginationRequest
	Name   string    `form:"name" binding:"omitempty"`
	Email  string    `form:"email" binding:"omitempty,email"`
	RoleID uuid.UUID `form:"roleId" binding:"omitempty"`
}

type UserListResponse PaginatedResponse[models.User]
