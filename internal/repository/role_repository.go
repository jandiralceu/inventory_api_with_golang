package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// RoleRepository defines the persistence contract for role-related data operations.
type RoleRepository interface {
	// Create persists a new role record in the database.
	Create(ctx context.Context, role *models.Role) (*models.Role, error)
	// Delete removes a role record from the database by its unique ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single role by its unique ID.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	// FindAll retrieves all roles currently defined in the system.
	FindAll(ctx context.Context) ([]models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository initializes a GORM-based implementation of RoleRepository.
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

const (
	roleIDQuery = "id = ?"
)

var _ RoleRepository = (*roleRepository)(nil)

// Create inserts a new role record into the database, mapping any database errors.
func (r *roleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return role, nil
}

// Delete removes a role record and returns [apperrors.ErrNotFound] if no record was deleted.
func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Role{}, roleIDQuery, id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}
	return nil
}

// FindByID retrieves a specific role using its unique UUID.
func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleIDQuery, id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &role, nil
}

// FindAll executes a query to retrieve all role records from the database.
func (r *roleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return roles, nil
}
