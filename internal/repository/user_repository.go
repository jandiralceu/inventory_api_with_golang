package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindAll(ctx context.Context, req dto.GetUserListRequest) (dto.PaginatedResponse[models.User], error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error
	ChangeRole(ctx context.Context, req dto.ChangeRoleRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

const (
	userIDQuery = "id = ?"
)

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, userIDQuery, id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userIDQuery, id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, req dto.GetUserListRequest) (dto.PaginatedResponse[models.User], error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+sanitizeLike(req.Name)+"%")
	}
	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+sanitizeLike(req.Email)+"%")
	}
	if req.RoleID != uuid.Nil {
		query = query.Where("role_id = ?", req.RoleID)
	}

	if err := query.Count(&total).Error; err != nil {
		return dto.PaginatedResponse[models.User]{}, mapDatabaseError(err)
	}

	page := req.GetPage()
	limit := req.GetLimit()
	offset := (page - 1) * limit
	order := req.GetSort("created_at", "name", "email") + " " + req.GetOrder()

	err := query.Preload("Role").
		Order(order).
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return dto.PaginatedResponse[models.User]{}, mapDatabaseError(err)
	}

	return dto.NewPaginatedResponse(users, total, page, limit), nil
}

func (r *userRepository) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	return errors.New("not implemented")
}

func (r *userRepository) ChangeRole(ctx context.Context, req dto.ChangeRoleRequest) error {
	return errors.New("not implemented")
}
