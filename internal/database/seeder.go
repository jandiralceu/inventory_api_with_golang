package database

import (
	"context"
	"log/slog"

	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// SeedRoles fills the database with initial roles if they don't exist.
func SeedRoles(ctx context.Context, db *gorm.DB) error {
	roles := []models.Role{
		{
			Name:        "admin",
			Description: "Full access: manages users, products, and system configurations",
		},
		{
			Name:        "manager",
			Description: "Manages products, suppliers, and views reports",
		},
		{
			Name:        "operator",
			Description: "Daily operations: records inventory movements and checks stock",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		// Check if role already exists by name
		err := db.WithContext(ctx).Where("name = ?", role.Name).First(&existingRole).Error

		if err == gorm.ErrRecordNotFound {
			slog.Info("Seeding role", "name", role.Name)
			if err := db.WithContext(ctx).Create(&role).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			slog.Debug("Role already exists, skipping", "name", role.Name)
		}
	}

	return nil
}
