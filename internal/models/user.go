package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system.
type User struct {
	// ID is the unique identifier for the user.
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Name is the full name of the user.
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	// Email is the unique email address of the user.
	Email string `gorm:"type:varchar(255);not null;unique" json:"email"`
	// PasswordHash is the argon2id encrypted password.
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	// RoleID is the foreign key linking the user to a role.
	RoleID uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	// Role is the hydrated role model representing the user's permissions.
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	// CreatedAt is the timestamp when the user account was created.
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	// UpdatedAt is the timestamp for the last update to the user profile.
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}
