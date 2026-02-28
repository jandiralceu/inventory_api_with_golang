package dto

import "github.com/google/uuid"

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type ChangeRoleRequest struct {
	UserID uuid.UUID `json:"userId" binding:"required"`
	RoleID uuid.UUID `json:"roleId" binding:"required"`
}

type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required,min=3,max=100"`
	Email    string    `json:"email" binding:"required,email,max=255"`
	Password string    `json:"password" binding:"required,min=8"`
	RoleID   uuid.UUID `json:"roleId" binding:"required"`
}

type GetUserListRequest struct {
	PaginationRequest
	Name   string    `form:"name" binding:"omitempty"`
	Email  string    `form:"email" binding:"omitempty,email"`
	RoleID uuid.UUID `form:"roleId" binding:"omitempty"`
}
