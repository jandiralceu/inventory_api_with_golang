package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines the persistence contract for user-related data operations.
type UserRepository interface {
	// Create persists a new user record in the database.
	Create(ctx context.Context, user *models.User) error
	// FindAll retrieves a list of users filtered by the provided criteria and supports pagination.
	FindAll(ctx context.Context, filter UserListFilter) (users []models.User, total int64, err error)
	// FindByID retrieves a single user along with their associated role by their unique ID.
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// FindByEmail locates a user and their role by their unique email address.
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// ChangePassword updates the password hash for the specified user record.
	ChangePassword(ctx context.Context, userID uuid.UUID, newHashedPassword string) error
	// ChangeRole updates the assigned role ID for a specific user.
	ChangeRole(ctx context.Context, userID uuid.UUID, newRoleID uuid.UUID) error
	// Delete removes a user record from the database using their unique ID.
	Delete(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

var _ UserRepository = (*userRepository)(nil)

// NewUserRepository initializes a GORM-based implementation of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// UserListFilter encapsulates search and pagination parameters for user listing operations.
type UserListFilter struct {
	Name       string
	Email      string
	RoleID     uuid.UUID
	Pagination PaginationParams
}

const (
	userIDQuery = "id = ?"
)

// Create inserts a new user record into the database, mapping any database errors.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Delete removes a user record and returns [apperrors.ErrNotFound] if no record was deleted.
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

// FindByID retrieves a user by ID and preloads the associated role.
func (r *userRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, userIDQuery, userID).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

// FindByEmail retrieves a user by their email address and preloads the associated role.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &user, nil
}

// FindAll executes a listing query with dynamic filtering, count calculation, and pagination.
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

// ChangePassword updates the password_hash field for the target user.
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

// ChangeRole updates the role_id field for the target user.
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
