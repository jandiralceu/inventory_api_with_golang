package dto

// CreateRoleRequest defines the payload for creating a new system role.
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=3"`
}
