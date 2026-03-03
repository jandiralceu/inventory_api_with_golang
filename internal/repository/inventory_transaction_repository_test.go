package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestInventoryTransactionRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryTransactionRepository(gormDB)

	id := uuid.New()
	inventoryID := uuid.New()
	productID := uuid.New()
	userID := uuid.New()
	warehouseID := uuid.New()

	transaction := &models.InventoryTransaction{
		ID:              id,
		InventoryID:     inventoryID,
		ProductID:       productID,
		WarehouseID:     warehouseID,
		UserID:          &userID,
		TransactionType: "ADJUSTMENT",
		QuantityChange:  10,
		QuantityBalance: 50,
		Reason:          "Correction",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "inventory_transactions"`)).
		WithArgs(
			transaction.InventoryID,
			transaction.ProductID,
			transaction.WarehouseID,
			transaction.UserID,
			transaction.QuantityChange,
			transaction.QuantityBalance,
			transaction.TransactionType,
			transaction.ReferenceID,
			transaction.Reason,
			transaction.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(transaction.ID, time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), nil, transaction)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryTransactionRepository_FindByInventoryID(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryTransactionRepository(gormDB)

	inventoryID := uuid.New()
	userID := uuid.New()
	params := PaginationParams{Limit: 10, Page: 1}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "inventory_transactions" WHERE inventory_id = $1`)).
		WithArgs(inventoryID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Match case-insensitive DESC
	mock.ExpectQuery(`SELECT \* FROM "inventory_transactions" WHERE inventory_id = \$1 ORDER BY created_at (?i)DESC LIMIT \$2`).
		WithArgs(inventoryID, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "inventory_id", "user_id"}).
			AddRow(uuid.New(), inventoryID, userID))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(userID, "Test User"))

	txs, total, err := repo.FindByInventoryID(context.Background(), inventoryID, params)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, txs, 1)
	assert.Equal(t, userID, *txs[0].UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryTransactionRepository_FindAll(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryTransactionRepository(gormDB)

	productID := uuid.New()
	userID := uuid.New()
	warehouseID := uuid.New()
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	filter := TransactionListFilter{
		ProductID:       &productID,
		UserID:          &userID,
		TransactionType: "INBOUND",
		StartDate:       &startDate,
		EndDate:         &endDate,
		Pagination:      PaginationParams{Limit: 10, Page: 1},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "inventory_transactions" WHERE product_id = $1 AND user_id = $2 AND transaction_type = $3 AND created_at >= $4 AND created_at <= $5`)).
		WithArgs(productID, userID, "INBOUND", startDate, endDate).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "inventory_transactions" WHERE product_id = \$1 AND user_id = \$2 AND transaction_type = \$3 AND created_at >= \$4 AND created_at <= \$5 ORDER BY created_at (?i)DESC LIMIT \$6`).
		WithArgs(productID, userID, "INBOUND", startDate, endDate, 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "user_id", "warehouse_id"}).
			AddRow(uuid.New(), productID, userID, warehouseID))

	// Preload Product
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."id" = $1`)).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(productID))

	// Preload User
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	// Preload Warehouse
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE "warehouses"."id" = $1`)).
		WithArgs(warehouseID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(warehouseID))

	txs, total, err := repo.FindAll(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, txs, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
