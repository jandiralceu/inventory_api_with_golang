package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// InventoryListFilter defines the criteria for filtering and paginating inventory records.
type InventoryListFilter struct {
	// ProductID filters inventory by a specific product unique ID.
	ProductID *uuid.UUID
	// WarehouseID filters inventory by a specific warehouse unique ID.
	WarehouseID *uuid.UUID
	// LowStock, if true, returns only items where quantity is less than or equal to min_quantity.
	LowStock *bool
	// Page is the current page number for pagination (starts at 1).
	Page int
	// Limit is the number of records per page.
	Limit int
	// Sort is the field name to use for sorting results.
	Sort string
	// Order is the sort direction (e.g., "asc" or "desc").
	Order string
}

// InventoryRepository defines the contract for inventory data persistence operations.
// All methods return mapped application errors via mapDatabaseError.
type InventoryRepository interface {
	// Create persists a new inventory record.
	// Returns ErrConflict if a record already exists for the product-warehouse pair.
	Create(ctx context.Context, inventory *models.Inventory) error
	// Update updates all fields of an existing inventory record.
	// Returns ErrNotFound if the record does not exist.
	Update(ctx context.Context, inventory *models.Inventory) error
	// Delete removes an inventory record by its unique ID.
	// Returns ErrNotFound if the record does not exist.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single inventory record with its associations preloaded.
	// Returns ErrNotFound if the record does not exist.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
	// FindAll returns a paginated list of inventory records matching the filter criteria.
	FindAll(ctx context.Context, filter InventoryListFilter) ([]models.Inventory, int64, error)
	// UpdateStock performs an atomic update on the quantity using optimistic locking.
	// Returns ErrConflict if the version has changed since the last read.
	UpdateStock(ctx context.Context, id uuid.UUID, quantityDelta int, version int) error
	// UpdateReservedStock performs an atomic update on the reserved quantity using optimistic locking.
	// Returns ErrConflict if the version has changed since the last read.
	UpdateReservedStock(ctx context.Context, id uuid.UUID, reservedDelta int, version int) error
}

// inventoryRepository is the GORM-based implementation of InventoryRepository.
type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new instance of the inventory repository.
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Create adds a new inventory entry to the database.
func (r *inventoryRepository) Create(ctx context.Context, inventory *models.Inventory) error {
	if err := r.db.WithContext(ctx).Create(inventory).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Update modifies an existing inventory record. It uses the ID to locate the record and
// updates all fields. Returns ErrNotFound if the record doesn't exist.
func (r *inventoryRepository) Update(ctx context.Context, inventory *models.Inventory) error {
	if err := r.db.WithContext(ctx).Save(inventory).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// Delete removes an inventory record by its unique identifier.
// Returns ErrNotFound if the record doesn't exist.
func (r *inventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Inventory{}, id).Error; err != nil {
		return mapDatabaseError(err)
	}
	return nil
}

// FindByID retrieves an inventory record by its ID, preloading Product and Warehouse relations.
func (r *inventoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.WithContext(ctx).Preload("Product").Preload("Warehouse").First(&inventory, id).Error
	if err != nil {
		return nil, mapDatabaseError(err)
	}
	return &inventory, nil
}

// FindAll retrieves a list of inventory records based on the provided filters.
func (r *inventoryRepository) FindAll(ctx context.Context, filter InventoryListFilter) ([]models.Inventory, int64, error) {
	var inventory []models.Inventory
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Inventory{})

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.LowStock != nil && *filter.LowStock {
		query = query.Where("quantity <= min_quantity")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	if filter.Sort != "" {
		order := "asc"
		if filter.Order == "desc" {
			order = "desc"
		}
		query = query.Order(filter.Sort + " " + order)
	} else {
		query = query.Order("updated_at desc")
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Preload("Product").Preload("Warehouse").Offset(offset).Limit(filter.Limit).Find(&inventory).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return inventory, total, nil
}

// UpdateStock increments or decrements the available quantity using optimistic locking.
// It matches the ID and the current version to ensure no concurrent updates occurred.
// Returns gorm.ErrRecordNotFound (mapped) if the version has changed.
func (r *inventoryRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantityDelta int, version int) error {
	result := r.db.WithContext(ctx).Model(&models.Inventory{}).
		Where("id = ? AND version = ?", id, version).
		Omit("updated_at").
		Updates(map[string]any{
			"quantity": gorm.Expr("quantity + ?", quantityDelta),
			"version":  gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound) // Concurrency conflict or not found
	}

	return nil
}

// UpdateReservedStock increments or decrements the reserved quantity using optimistic locking.
// Returns gorm.ErrRecordNotFound if the version does not match, indicating a concurrent update.
func (r *inventoryRepository) UpdateReservedStock(ctx context.Context, id uuid.UUID, reservedDelta int, version int) error {
	result := r.db.WithContext(ctx).Model(&models.Inventory{}).
		Where("id = ? AND version = ?", id, version).
		Omit("updated_at").
		Updates(map[string]any{
			"reserved_quantity": gorm.Expr("reserved_quantity + ?", reservedDelta),
			"version":           gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return mapDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return mapDatabaseError(gorm.ErrRecordNotFound)
	}

	return nil
}
