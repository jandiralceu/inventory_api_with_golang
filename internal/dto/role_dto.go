package dto

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=3"`
}
