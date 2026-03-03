package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"gorm.io/gorm"
)

// InventoryTransactionRepository defines the contract for auditing stock movements.
type InventoryTransactionRepository interface {
	// Create records a new stock movement in the audit log.
	// Can be used within a transaction.
	Create(ctx context.Context, tx *gorm.DB, transaction *models.InventoryTransaction) error

	// FindByInventoryID retrieves all transactions for a specific inventory record.
	// Supports pagination and filtering.
	FindByInventoryID(ctx context.Context, inventoryID uuid.UUID, params PaginationParams) ([]models.InventoryTransaction, int64, error)

	// FindAll retrieves a global list of transactions with advanced filtering.
	FindAll(ctx context.Context, filter TransactionListFilter) ([]models.InventoryTransaction, int64, error)
}

// TransactionListFilter encapsulates filters for the global audit log.
type TransactionListFilter struct {
	InventoryID     *uuid.UUID
	ProductID       *uuid.UUID
	WarehouseID     *uuid.UUID
	UserID          *uuid.UUID
	TransactionType string
	StartDate       *time.Time
	EndDate         *time.Time
	Pagination      PaginationParams
}

type inventoryTransactionRepository struct {
	db *gorm.DB
}

// NewInventoryTransactionRepository creates a new instance of the audit repository.
func NewInventoryTransactionRepository(db *gorm.DB) InventoryTransactionRepository {
	return &inventoryTransactionRepository{db: db}
}

// Create records a new stock movement in the audit log. It accepts an optional *gorm.DB
// to allow participation in an existing transaction.
func (r *inventoryTransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *models.InventoryTransaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(transaction).Error; err != nil {
		return mapDatabaseError(err)
	}

	return nil
}

// FindByInventoryID retrieves all transactions for a specific inventory record with pagination.
func (r *inventoryTransactionRepository) FindByInventoryID(ctx context.Context, inventoryID uuid.UUID, params PaginationParams) ([]models.InventoryTransaction, int64, error) {
	var transactions []models.InventoryTransaction
	var total int64

	db := r.db.WithContext(ctx).Model(&models.InventoryTransaction{}).Where("inventory_id = ?", inventoryID)

	// Count total records for pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	// Apply pagination and latest first
	query := db.Order(params.GetOrderBy()).
		Limit(params.Limit).
		Offset(params.GetOffset())

	// Optional: Preload User and other relations if needed for the UI
	if err := query.Preload("User").Find(&transactions).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return transactions, total, nil
}

// FindAll retrieves a global list of transactions with advanced filtering, preloads, and pagination.
func (r *inventoryTransactionRepository) FindAll(ctx context.Context, filter TransactionListFilter) ([]models.InventoryTransaction, int64, error) {
	var transactions []models.InventoryTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.InventoryTransaction{})

	// Apply Filters
	if filter.InventoryID != nil {
		query = query.Where("inventory_id = ?", *filter.InventoryID)
	}
	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.TransactionType != "" {
		query = query.Where("transaction_type = ?", filter.TransactionType)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	// Apply pagination and preloads
	err := query.Order(filter.Pagination.GetOrderBy()).
		Limit(filter.Pagination.Limit).
		Offset(filter.Pagination.GetOffset()).
		Preload("User").
		Preload("Product").
		Preload("Warehouse").
		Find(&transactions).Error

	if err != nil {
		return nil, 0, mapDatabaseError(err)
	}

	return transactions, total, nil
}
