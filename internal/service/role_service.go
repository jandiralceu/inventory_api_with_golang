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

type RoleService interface {
	Create(ctx context.Context, role *models.Role) (*models.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	cache    pkg.CacheManager
}

func NewRoleService(roleRepo repository.RoleRepository, cache pkg.CacheManager) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		cache:    cache,
	}
}

const (
	roleCachePrefix = "role:"
)

var _ RoleService = (*roleService)(nil)

func (s *roleService) Create(ctx context.Context, role *models.Role) (*models.Role, error) {
	create, err := s.roleRepo.Create(ctx, role)
	if err == nil {
		s.cache.DeletePrefix(ctx, roleCachePrefix)
	}

	return create, err
}

func (s *roleService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.roleRepo.Delete(ctx, id)
	if err == nil {
		s.cache.DeletePrefix(ctx, roleCachePrefix)
	}

	return err
}

func (s *roleService) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	cacheKey := fmt.Sprintf("%sid:%s", roleCachePrefix, id)
	var cached models.Role
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return &cached, nil
	}

	role, err := s.roleRepo.FindByID(ctx, id)
	if err == nil {
		s.cache.Set(ctx, cacheKey, role, 72*time.Hour)
	}

	return role, err
}

func (s *roleService) FindAll(ctx context.Context) ([]models.Role, error) {
	cacheKey := fmt.Sprintf("%slist", roleCachePrefix)
	var cached []models.Role
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		return cached, nil
	}

	roles, err := s.roleRepo.FindAll(ctx)
	if err == nil {
		s.cache.Set(ctx, cacheKey, roles, 72*time.Hour)
	}

	return roles, err
}
