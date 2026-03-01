package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	user := &models.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "plain-password",
	}

	mockHasher.On("Hash", "plain-password").Return("hashed-password", nil)
	mockRepo.On("Create", ctx, user).Return(nil)

	err := svc.Create(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, "hashed-password", user.PasswordHash)
	mockHasher.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestFindAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	req := dto.GetUserListRequest{
		Name: "test",
	}

	users := []models.User{{Name: "test user"}}
	mockRepo.On("FindAll", ctx, mock.AnythingOfType("repository.UserListFilter")).
		Return(users, int64(1), nil)

	res, err := svc.FindAll(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Len(t, res.Data, 1)
	mockRepo.AssertExpectations(t)
}

func TestFindUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	userID := uuid.New()
	user := &models.User{ID: userID, Name: "test user"}

	mockRepo.On("FindByID", ctx, userID).Return(user, nil)

	res, err := svc.FindByID(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := svc.Delete(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestChangePassword(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, mockHasher)

		ctx := context.Background()
		userID := uuid.New()
		req := dto.ChangePasswordRequest{
			OldPassword: "old",
			NewPassword: "new",
		}

		user := &models.User{ID: userID, PasswordHash: "old-hashed"}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)
		mockHasher.On("Verify", "old", "old-hashed").Return(true, nil)
		mockHasher.On("Hash", "new").Return("new-hashed", nil)
		mockRepo.On("ChangePassword", ctx, userID, "new-hashed").Return(nil)

		err := svc.ChangePassword(ctx, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("IncorrectOldPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		mockHasher := new(MockPasswordHasher)
		svc := NewUserService(mockRepo, mockHasher)

		ctx := context.Background()
		userID := uuid.New()
		req := dto.ChangePasswordRequest{
			OldPassword: "wrong",
			NewPassword: "new",
		}

		user := &models.User{ID: userID, PasswordHash: "old-hashed"}

		mockRepo.On("FindByID", ctx, userID).Return(user, nil)
		mockHasher.On("Verify", "wrong", "old-hashed").Return(false, nil)

		err := svc.ChangePassword(ctx, userID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperrors.ErrUnauthorized))
	})
}

func TestChangeRole(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)
	svc := NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	userID := uuid.New()
	roleID := uuid.New()
	req := dto.ChangeRoleRequest{
		UserID: userID,
		RoleID: roleID,
	}

	mockRepo.On("ChangeRole", ctx, userID, roleID).Return(nil)

	err := svc.ChangeRole(ctx, userID, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
