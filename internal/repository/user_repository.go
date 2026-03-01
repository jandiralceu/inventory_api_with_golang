package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindAll(ctx context.Context, filter UserListFilter) (users []models.User, total int64, err error)
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error
	ChangeRole(ctx context.Context, userID uuid.UUID, newRoleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type UserListFilter struct {
	Name       string
	Email      string
	RoleID     uuid.UUID
	Pagination PaginationParams
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

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, userIDQuery, userID)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, userIDQuery, userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, filter UserListFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+sanitizeLike(filter.Name)+"%")
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+sanitizeLike(filter.Email)+"%")
	}
	if filter.RoleID != uuid.Nil {
		query = query.Where("role_id = ?", filter.RoleID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	err := query.Preload("Role").
		Order(filter.Pagination.GetOrderBy()).
		Offset(filter.Pagination.GetOffset()).
		Limit(filter.Pagination.Limit).
		Find(&users).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return users, total, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where(userIDQuery, userID).
		Update("password_hash", newHashedPassword)

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *userRepository) ChangeRole(ctx context.Context, userID uuid.UUID, newRoleID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where(userIDQuery, userID).
		Update("role_id", newRoleID)

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
