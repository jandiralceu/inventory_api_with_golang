package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

// UserService defines the business logic contract for user management.
type UserService interface {
	// Create registers a new user, ensuring the password is securely hashed.
	Create(ctx context.Context, user *models.User) error
	// FindAll returns a paginated list of users based on search criteria.
	FindAll(ctx context.Context, req dto.GetUserListRequest) (dto.PaginatedResponse[models.User], error)
	// FindByID retrieves a single user by their unique UUID.
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// FindByEmail locates a user using their unique email address.
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	// ChangePassword updates a user's password after verifying the old one.
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error
	// ChangeRole updates a user's assigned role.
	ChangeRole(ctx context.Context, userID uuid.UUID, req dto.ChangeRoleRequest) error
	// Delete removes a user from the system.
	Delete(ctx context.Context, userID uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
	hasher   pkg.PasswordHasher
}

var _ UserService = (*userService)(nil)

// NewUserService initializes a UserService with its required repository and hasher dependencies.
func NewUserService(userRepo repository.UserRepository, hasher pkg.PasswordHasher) UserService {
	return &userService{userRepo: userRepo, hasher: hasher}
}

// Create performs password hashing using the injected hasher before persisting the user through the repository.
func (s *userService) Create(ctx context.Context, user *models.User) error {
	hashedPassword, err := s.hasher.Hash(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword

	return s.userRepo.Create(ctx, user)
}

// FindAll delegates the retrieval of the user list to the repository.
func (s *userService) FindAll(ctx context.Context, req dto.GetUserListRequest) (dto.PaginatedResponse[models.User], error) {
	filter := repository.UserListFilter{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleID,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at", "name", "email"),
			Order: req.GetOrder(),
		},
	}

	users, total, err := s.userRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[models.User]{}, err
	}

	return dto.NewPaginatedResponse(users, total, filter.Pagination.Page, filter.Pagination.Limit), nil
}

// FindByID retrieves a user by their primary key from the repository.
func (s *userService) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// FindByEmail retrieves a user by their email address for authentication or identification purposes.
func (s *userService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// Delete removes the user record identified by the unique ID.
func (s *userService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	match, err := s.hasher.Verify(req.OldPassword, user.PasswordHash)
	if err != nil {
		return err
	}

	if !match {
		return apperrors.ErrUnauthorized
	}

	newHashedPassword, err := s.hasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.ChangePassword(ctx, userID, newHashedPassword)
}

func (s *userService) ChangeRole(ctx context.Context, userID uuid.UUID, req dto.ChangeRoleRequest) error {
	return s.userRepo.ChangeRole(ctx, userID, req.RoleID)
}
