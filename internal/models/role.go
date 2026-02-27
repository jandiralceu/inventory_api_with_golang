package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a user role in the system.
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
}

// BeforeSave normalizes the role name to lowercase before persisting.
func (r *Role) BeforeSave(tx *gorm.DB) error {
	r.Name = strings.ToLower(r.Name)
	return nil
}
