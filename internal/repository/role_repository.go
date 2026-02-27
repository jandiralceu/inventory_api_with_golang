package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

const (
	roleIDQuery = "id = ?"
)

var _ RoleRepository = (*roleRepository)(nil)

func (r *roleRepository) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Role{}, roleIDQuery, id)
	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *roleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleIDQuery, id).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, mapDatabaseError(err)
	}
	return roles, nil
}
