package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock UserService ---

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) FindAll(ctx context.Context, req dto.GetUserListRequest) (dto.PaginatedResponse[models.User], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.PaginatedResponse[models.User]), args.Error(1)
}

func (m *MockUserService) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) ChangeRole(ctx context.Context, userID uuid.UUID, req dto.ChangeRoleRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// --- Mock PasswordHasher ---

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(password, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}

// --- Mock CacheManager ---

type MockCacheManager struct {
	mock.Mock
}

func (m *MockCacheManager) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheManager) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheManager) DeletePrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}

func (m *MockCacheManager) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCacheManager) GetClient() *redis.Client {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.Client)
}

// --- Test helpers ---

func generateTestRSAKeys(t *testing.T) (*pkg.JWTManager, string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	jwtManager, err := pkg.NewJWTManager(string(privPEM), string(pubPEM))
	if err != nil {
		t.Fatalf("failed to create JWT manager: %v", err)
	}

	return jwtManager, string(privPEM), string(pubPEM)
}

func setupAuthHandler(t *testing.T) (*AuthHandler, *MockUserService, *MockCacheManager, *pkg.JWTManager) {
	t.Helper()
	mockService := new(MockUserService)
	mockCache := new(MockCacheManager)
	jwtManager, _, _ := generateTestRSAKeys(t)
	hasher := pkg.NewArgon2PasswordHasher()
	handler := NewAuthHandler(mockService, jwtManager, mockCache, hasher)
	return handler, mockService, mockCache, jwtManager
}

func setupAuthHandlerWithMockHasher(t *testing.T) (*AuthHandler, *MockUserService, *MockCacheManager, *MockPasswordHasher, *pkg.JWTManager) {
	t.Helper()
	mockService := new(MockUserService)
	mockCache := new(MockCacheManager)
	mockHasher := new(MockPasswordHasher)
	jwtManager, _, _ := generateTestRSAKeys(t)
	handler := NewAuthHandler(mockService, jwtManager, mockCache, mockHasher)
	return handler, mockService, mockCache, mockHasher, jwtManager
}

// =====================
// Register Tests
// =====================

func TestRegister_Success(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	roleID := uuid.New()
	mockService.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, apperrors.ErrNotFound)
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	body := map[string]any{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
		"roleId":   roleID.String(),
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestRegister_BadRequest_MissingFields(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	// Missing "password" and "roleId"
	body := map[string]any{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestRegister_BadRequest_InvalidEmail(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	body := map[string]any{
		"name":     "John Doe",
		"email":    "not-an-email",
		"password": "password123",
		"roleId":   uuid.New().String(),
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestRegister_BadRequest_PasswordTooShort(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	body := map[string]any{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "short",
		"roleId":   uuid.New().String(),
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Create")
}

func TestRegister_Conflict_EmailExists(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	mockService.On("FindByEmail", mock.Anything, "john@example.com").
		Return(&models.User{}, nil)

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	body := map[string]any{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
		"roleId":   uuid.New().String(),
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Conflict", resp.Title)
	assert.Equal(t, "resource conflict: email already in use", resp.Detail)
	mockService.AssertNotCalled(t, "Create")
}

func TestRegister_InternalServerError(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	mockService.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, apperrors.ErrNotFound)
	mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
		Return(errors.New("database error"))

	router := setupRouter()
	router.POST("/auth/register", handler.Register)

	body := map[string]any{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
		"roleId":   uuid.New().String(),
	}

	w := performRequest(router, "POST", "/auth/register", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Internal Server Error", resp.Title)
	assert.Equal(t, "An unexpected error occurred. Please try again later.", resp.Detail)
	mockService.AssertExpectations(t)
}

// =====================
// SignIn Tests
// =====================

func TestSignIn_Success(t *testing.T) {
	handler, mockService, mockCache, _ := setupAuthHandler(t)

	hasher := pkg.NewArgon2PasswordHasher()
	hashedPassword, err := hasher.Hash("password123")
	assert.NoError(t, err)

	userID := uuid.New()
	mockService.On("FindByEmail", mock.Anything, "john@example.com").
		Return(&models.User{
			ID:           userID,
			Name:         "John Doe",
			Email:        "john@example.com",
			PasswordHash: hashedPassword,
			Role:         models.Role{Name: "admin"},
		}, nil)

	mockCache.On("Set", mock.Anything, mock.Anything, "active", mock.Anything).Return(nil)

	router := setupRouter()
	router.POST("/auth/signin", handler.SignIn)

	body := map[string]any{
		"email":    "john@example.com",
		"password": "password123",
	}

	w := performRequest(router, "POST", "/auth/signin", body)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["accessToken"])
	assert.NotEmpty(t, response["refreshToken"])
	mockService.AssertExpectations(t)
}

func TestSignIn_BadRequest_MissingFields(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/signin", handler.SignIn)

	body := map[string]any{
		"email": "john@example.com",
	}

	w := performRequest(router, "POST", "/auth/signin", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "FindByEmail")
}

func TestSignIn_Unauthorized_UserNotFound(t *testing.T) {
	handler, mockService, _, _, _ := setupAuthHandlerWithMockHasher(t)

	mockService.On("FindByEmail", mock.Anything, "unknown@example.com").
		Return(nil, errors.New("user not found"))

	router := setupRouter()
	router.POST("/auth/signin", handler.SignIn)

	body := map[string]any{
		"email":    "unknown@example.com",
		"password": "password123",
	}

	w := performRequest(router, "POST", "/auth/signin", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Unauthorized", resp.Title)
	assert.Equal(t, "unauthorized: invalid email or password", resp.Detail)
	mockService.AssertExpectations(t)
}

func TestSignIn_Unauthorized_WrongPassword(t *testing.T) {
	handler, mockService, _, _ := setupAuthHandler(t)

	hasher := pkg.NewArgon2PasswordHasher()
	hashedPassword, _ := hasher.Hash("correct-password")

	mockService.On("FindByEmail", mock.Anything, "john@example.com").
		Return(&models.User{
			ID:           uuid.New(),
			Email:        "john@example.com",
			PasswordHash: hashedPassword,
		}, nil)

	router := setupRouter()
	router.POST("/auth/signin", handler.SignIn)

	body := map[string]any{
		"email":    "john@example.com",
		"password": "wrong-password",
	}

	w := performRequest(router, "POST", "/auth/signin", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Unauthorized", resp.Title)
	assert.Equal(t, "unauthorized: invalid email or password", resp.Detail)
	mockService.AssertExpectations(t)
}

// =====================
// SignOut Tests
// =====================

func TestSignOut_Success(t *testing.T) {
	handler, _, mockCache, _, jwtManager := setupAuthHandlerWithMockHasher(t)

	userID := uuid.New()
	validRefreshToken, _ := jwtManager.GenerateToken(userID, "admin", _refreshTokenExpiration)
	refreshKey := fmt.Sprintf("refresh_token:%s:%s", userID.String(), validRefreshToken)

	mockCache.On("Delete", mock.Anything, refreshKey).Return(nil)

	router := setupRouter()
	router.POST("/auth/signout", handler.SignOut)

	body := map[string]any{
		"refreshToken": validRefreshToken,
	}

	w := performRequest(router, "POST", "/auth/signout", body)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockCache.AssertExpectations(t)
}

func TestSignOut_BadRequest_MissingToken(t *testing.T) {
	handler, _, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/signout", handler.SignOut)

	body := map[string]any{}

	w := performRequest(router, "POST", "/auth/signout", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignOut_Unauthorized_InvalidToken(t *testing.T) {
	handler, _, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/signout", handler.SignOut)

	body := map[string]any{
		"refreshToken": "invalid-token",
	}

	w := performRequest(router, "POST", "/auth/signout", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Unauthorized", resp.Title)
	assert.Equal(t, "unauthorized: invalid or expired refresh token", resp.Detail)
}

// =====================
// RefreshToken Tests
// =====================

func TestRefreshToken_Success(t *testing.T) {
	handler, mockService, mockCache, _, jwtManager := setupAuthHandlerWithMockHasher(t)

	userID := uuid.New()

	validRefreshToken, err := jwtManager.GenerateToken(userID, "admin", _refreshTokenExpiration)
	assert.NoError(t, err)

	mockService.On("FindByID", mock.Anything, userID).
		Return(&models.User{ID: userID, Name: "John Doe", Role: models.Role{Name: "admin"}}, nil)

	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)
	mockCache.On("Set", mock.Anything, mock.Anything, "active", mock.Anything).Return(nil)

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{
		"refreshToken": validRefreshToken,
	}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["accessToken"])
	assert.NotEmpty(t, response["refreshToken"])
	mockService.AssertExpectations(t)
}

func TestRefreshToken_BadRequest_MissingToken(t *testing.T) {
	handler, _, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefreshToken_Unauthorized_InvalidToken(t *testing.T) {
	handler, _, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{
		"refreshToken": "invalid.token.value",
	}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Unauthorized", resp.Title)
	assert.Equal(t, "unauthorized: invalid or expired refresh token", resp.Detail)
}

func TestRefreshToken_Unauthorized_TokenNotInCache(t *testing.T) {
	handler, _, mockCache, _, jwtManager := setupAuthHandlerWithMockHasher(t)

	userID := uuid.New()

	validRefreshToken, err := jwtManager.GenerateToken(userID, "admin", _refreshTokenExpiration)
	assert.NoError(t, err)

	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("key not found"))

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{
		"refreshToken": validRefreshToken,
	}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "unauthorized: refresh token not found or already used", resp.Detail)
}

func TestRefreshToken_Unauthorized_UserNotFound(t *testing.T) {
	handler, mockService, mockCache, _, jwtManager := setupAuthHandlerWithMockHasher(t)

	userID := uuid.New()

	validRefreshToken, err := jwtManager.GenerateToken(userID, "admin", _refreshTokenExpiration)
	assert.NoError(t, err)

	mockService.On("FindByID", mock.Anything, userID).
		Return(nil, errors.New("user not found"))

	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{
		"refreshToken": validRefreshToken,
	}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "unauthorized: user no longer exists", resp.Detail)
	mockService.AssertExpectations(t)
}

func TestRefreshToken_InternalServerError_FailedToSaveNewToken(t *testing.T) {
	handler, mockService, mockCache, _, jwtManager := setupAuthHandlerWithMockHasher(t)

	userID := uuid.New()

	validRefreshToken, err := jwtManager.GenerateToken(userID, "admin", _refreshTokenExpiration)
	assert.NoError(t, err)

	mockService.On("FindByID", mock.Anything, userID).
		Return(&models.User{ID: userID, Name: "John Doe", Role: models.Role{Name: "admin"}}, nil)

	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockCache.On("Delete", mock.Anything, mock.Anything).Return(nil)
	mockCache.On("Set", mock.Anything, mock.Anything, "active", mock.Anything).
		Return(errors.New("redis error"))

	router := setupRouter()
	router.POST("/auth/refresh", handler.RefreshToken)

	body := map[string]any{
		"refreshToken": validRefreshToken,
	}

	w := performRequest(router, "POST", "/auth/refresh", body)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp ProblemDetails
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Internal Server Error", resp.Title)
	assert.Equal(t, "An unexpected error occurred. Please try again later.", resp.Detail)
}

// =====================
// Gin Route Registration
// =====================

func TestAuthRoutes_MethodNotAllowed(t *testing.T) {
	handler, _, _, _, _ := setupAuthHandlerWithMockHasher(t)

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/signin", handler.SignIn)

	// GET on a POST-only route
	w := performRequest(router, "GET", "/auth/register", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
