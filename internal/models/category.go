package models

import (
	"time"

	"github.com/google/uuid"
)

// Category reflects the organizational hierarchy for products.
// It supports recursive structures through the Parent relationship.
type Category struct {
	// ID is the unique identifier for the category.
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Name is the display name of the category (e.g., "Electronics").
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
	// Slug is the URL-friendly name used for queries (e.g., "electronics").
	Slug string `gorm:"type:varchar(100);not null;unique" json:"slug"`
	// Description provides additional context about the category.
	Description string `gorm:"type:text" json:"description"`
	// ParentID links to the parent category for hierarchical grouping.
	ParentID *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	// IsActive determines if the category is visible in the frontend.
	IsActive bool `gorm:"type:boolean;not null;default:true" json:"is_active"`
	// CreatedAt is the timestamp when the category was persistent.
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	// UpdatedAt is the timestamp for the last update to the category record.
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}
