package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Email        string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}
