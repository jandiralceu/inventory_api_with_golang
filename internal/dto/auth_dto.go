package dto

import "github.com/jandiralceu/inventory_api_with_golang/internal/models"

// RegisterRequest defines the input data for new user registration.
type RegisterRequest struct {
	CreateUserRequest
}

// RegisterResponse represents the response sent back after a successful user registration.
type RegisterResponse struct {
	models.User
}

// SignInRequest defines the credentials used for authentication.
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignInResponse contains the tokens issued upon successful authentication.
type SignInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenRequest defines the payload for rotation of expired access tokens.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse contains the new token pair issued after refresh.
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// SignOutRequest defines the payload to revoke a session via refresh token.
type SignOutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
