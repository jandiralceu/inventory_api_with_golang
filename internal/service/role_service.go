package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	pkg "github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

// RoleService defines the business logic contract for role management, including cache orchestration.
type RoleService interface {
	// Create persists a new role and invalidates the role cache.
	Create(ctx context.Context, role *models.Role) (*models.Role, error)
	// Delete removes a role by ID and invalidates the role cache.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a role by ID, checking the cache before the repository.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	// FindAll retrieves all roles, checking the cache before the repository.
	FindAll(ctx context.Context) ([]models.Role, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	cache    pkg.CacheManager
}

// NewRoleService initializes a RoleService with repository and cache dependencies.
func NewRoleService(roleRepo repository.RoleRepository, cache pkg.CacheManager) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		cache:    cache,
	}
}

const (
	_roleCachePrefix = "role:"
)

var _ RoleService = (*roleService)(nil)

// Create delegates role creation to the repository and purges the role cache on success.
func (s *roleService) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	create, err := s.roleRepo.Create(ctx, role)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, _roleCachePrefix)
	}

	return create, err
}

// Delete removes a role via the repository and purges the role cache on success.
func (s *roleService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.roleRepo.Delete(ctx, id)
	if err == nil {
		_ = s.cache.DeletePrefix(ctx, _roleCachePrefix)
	}

	return err
}

// FindByID retrieves a role from cache or falls back to the repository, populating the cache on success.
func (s *roleService) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	cacheKey := fmt.Sprintf("%sid:%s", _roleCachePrefix, id)
	var cached models.Role
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	role, err := s.roleRepo.FindByID(ctx, id)
	if err == nil {
		_ = s.cache.Set(ctx, cacheKey, role, 72*time.Hour)
	}

	return role, err
}

// FindAll retrieves all roles from cache or falls back to the repository, populating the cache on success.
func (s *roleService) FindAll(ctx context.Context) ([]models.Role, error) {
	cacheKey := fmt.Sprintf("%slist", _roleCachePrefix)
	var cached []models.Role
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	roles, err := s.roleRepo.FindAll(ctx)
	if err == nil {
		_ = s.cache.Set(ctx, cacheKey, roles, 72*time.Hour)
	}

	return roles, err
}
