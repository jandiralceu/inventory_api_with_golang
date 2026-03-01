package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/service"
)

const (
	_accessTokenExpiration  = 15 * time.Minute
	_refreshTokenExpiration = 7 * 24 * time.Hour
	_refreshTokenCacheKey   = "refresh_token:"
)

// AuthHandler manages identity operations including registration, login, and session rotation.
type AuthHandler struct {
	userService service.UserService
	jwtManager  *pkg.JWTManager
	cache       pkg.CacheManager
	hasher      pkg.PasswordHasher
}

// NewAuthHandler initializes an AuthHandler with its required dependencies for user management, JWT, cache, and hashing.
func NewAuthHandler(userService service.UserService, jwtManager *pkg.JWTManager, cache pkg.CacheManager, hasher pkg.PasswordHasher) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
		cache:       cache,
		hasher:      hasher,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with name, email, password and assigned role.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User registration data"
// @Success      204 "User registered successfully"
// @Failure      400 {object} ProblemDetails "Invalid input or validation error"
// @Failure      409 {object} ProblemDetails "Email already in use"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	// Verify if user already exists
	existingUser, _ := h.userService.FindByEmail(c.Request.Context(), req.Email)
	if existingUser != nil {
		RespondWithError(c, fmt.Errorf("%w: email already in use", apperrors.ErrConflict))
		return
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		RoleID:       req.RoleID,
	}

	if err := h.userService.Create(c.Request.Context(), user); err != nil {
		RespondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// SignIn godoc
// @Summary      Authenticate user
// @Description  Logs in a user with email and password, returning access and refresh tokens.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignInRequest true "Login credentials"
// @Success      200 {object} dto.SignInResponse "Tokens generated successfully"
// @Failure      400 {object} ProblemDetails "Invalid request format"
// @Failure      401 {object} ProblemDetails "Invalid email or password"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal server error"
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	// Find the user by email.
	user, err := h.userService.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: invalid email or password", apperrors.ErrUnauthorized))
		return
	}

	// Verify the password against the stored hash.
	match, err := h.hasher.Verify(req.Password, user.PasswordHash)
	if err != nil || !match {
		RespondWithError(c, fmt.Errorf("%w: invalid email or password", apperrors.ErrUnauthorized))
		return
	}

	// Generate access and refresh tokens.
	accessToken, err := h.jwtManager.GenerateToken(user.ID, user.Role.Name, _accessTokenExpiration)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to generate access token", apperrors.ErrInternal))
		return
	}

	refreshToken, err := h.jwtManager.GenerateToken(user.ID, user.Role.Name, _refreshTokenExpiration)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to generate refresh token", apperrors.ErrInternal))
		return
	}

	// Save the refresh token to Redis.
	refreshKey := fmt.Sprintf("%s%s:%s", _refreshTokenCacheKey, user.ID.String(), refreshToken)
	if err := h.cache.Set(c.Request.Context(), refreshKey, "active", _refreshTokenExpiration); err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to save refresh token", apperrors.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// SignOut godoc
// @Summary      User logout
// @Description  Invalidates the session by deleting the refresh token from the server-side cache.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignOutRequest true "Token to invalidate"
// @Success      204 "Successfully logged out"
// @Failure      400 {object} ProblemDetails "Invalid request format"
// @Failure      401 {object} ProblemDetails "Invalid or expired refresh token"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Router       /auth/signout [post]
func (h *AuthHandler) SignOut(c *gin.Context) {
	var req dto.SignOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: invalid or expired refresh token", apperrors.ErrUnauthorized))
		return
	}

	refreshKey := fmt.Sprintf("%s%s:%s", _refreshTokenCacheKey, claims.UserID.String(), req.RefreshToken)
	_ = h.cache.Delete(c.Request.Context(), refreshKey)

	c.Status(http.StatusNoContent)
}

// RefreshToken godoc
// @Summary      Refresh session
// @Description  Rotates the session by issuing a new access and refresh token pair using a valid refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} dto.RefreshTokenResponse "New tokens generated successfully"
// @Failure      400 {object} ProblemDetails "Invalid request format"
// @Failure      401 {object} ProblemDetails "Invalid, expired or already used refresh token"
// @Failure      429 {object} ProblemDetails "Too many requests"
// @Failure      500 {object} ProblemDetails "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, ParseValidationError(err))
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: invalid or expired refresh token", apperrors.ErrUnauthorized))
		return
	}

	refreshKey := fmt.Sprintf("%s%s:%s", _refreshTokenCacheKey, claims.UserID.String(), req.RefreshToken)
	var cachedToken string
	err = h.cache.Get(c.Request.Context(), refreshKey, &cachedToken)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: refresh token not found or already used", apperrors.ErrUnauthorized))
		return
	}

	// Verify the user still exists mapping from the id
	user, err := h.userService.FindByID(c.Request.Context(), claims.UserID)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: user no longer exists", apperrors.ErrUnauthorized))
		return
	}

	accessToken, err := h.jwtManager.GenerateToken(claims.UserID, user.Role.Name, _accessTokenExpiration)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to generate access token", apperrors.ErrInternal))
		return
	}

	refreshToken, err := h.jwtManager.GenerateToken(claims.UserID, user.Role.Name, _refreshTokenExpiration)
	if err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to generate refresh token", apperrors.ErrInternal))
		return
	}

	// Invalidate the old refresh token
	_ = h.cache.Delete(c.Request.Context(), refreshKey)

	// Save the new refresh token
	newRefreshKey := fmt.Sprintf("%s%s:%s", _refreshTokenCacheKey, claims.UserID.String(), refreshToken)
	if err := h.cache.Set(c.Request.Context(), newRefreshKey, "active", _refreshTokenExpiration); err != nil {
		RespondWithError(c, fmt.Errorf("%w: failed to save new refresh token", apperrors.ErrInternal))
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
