package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Address represents a physical location or contact coordinates for a supplier.
// It is stored as a JSONB column in the database.
type Address struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	ZipCode    string `json:"zip_code"`
}

// Value implements the driver.Valuer interface for database storage.
// This allows GORM to save the struct as a JSON string in the database.
func (a Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database retrieval.
// This allows GORM to unmarshal the JSON string from the database into the struct.
func (a *Address) Scan(value any) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Supplier represents a company or individual that provides products to the inventory.
type Supplier struct {
	// ID is the unique identifier for the supplier.
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	// Name is the display name of the supplier.
	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
	// Slug is the URL-friendly name used for queries.
	Slug string `gorm:"type:varchar(100);not null;unique" json:"slug"`
	// Description provides additional context about the supplier.
	Description string `gorm:"type:text" json:"description"`
	// TaxID is the fiscal identification (e.g., CNPJ, CPF, EIN).
	TaxID string `gorm:"type:varchar(50);unique" json:"tax_id"`
	// Email is the main contact email for the supplier.
	Email string `gorm:"type:varchar(255)" json:"email"`
	// Phone is the primary contact phone number.
	Phone string `gorm:"type:varchar(20)" json:"phone"`
	// Address is a flexible JSONB field containing location details.
	Address Address `gorm:"type:jsonb" json:"address"`
	// ContactPerson is the name of the primary representative at the supplier.
	ContactPerson string `gorm:"type:varchar(100)" json:"contact_person"`
	// IsActive determines if the supplier is currently enabled for operations.
	IsActive bool `gorm:"type:boolean;not null;default:true" json:"is_active"`
	// CreatedAt is the timestamp when the supplier was registered.
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	// UpdatedAt is the timestamp for the last update to the supplier record.
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`
}
