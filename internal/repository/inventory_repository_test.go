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

func TestInventoryRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()
	inventory := &models.Inventory{
		ID:          id,
		ProductID:   productID,
		WarehouseID: warehouseID,
		Quantity:    100,
		MinQuantity: 10,
		Version:     1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "inventory"`)).
		WithArgs(inventory.ProductID, inventory.WarehouseID, inventory.Quantity, inventory.ReservedQuantity, inventory.LocationCode, inventory.MinQuantity, inventory.MaxQuantity, inventory.Version, inventory.LastCountedAt, inventory.Metadata, sqlmock.AnyArg(), sqlmock.AnyArg(), inventory.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(inventory.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), inventory)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepository_Update(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()
	inventory := &models.Inventory{
		ID:           id,
		LocationCode: "NEW-LOC",
		UpdatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "inventory" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), inventory)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepository_Delete(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "inventory" WHERE "inventory"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepository_FindByID(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()
	productID := uuid.New()
	warehouseID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "inventory" WHERE "inventory"."id" = $1`)).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "product_id", "warehouse_id"}).
			AddRow(id, productID, warehouseID))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."id" = $1`)).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(productID, "Test Product"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE "warehouses"."id" = $1`)).
		WithArgs(warehouseID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(warehouseID, "Main Warehouse"))

	inventory, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, inventory)
	assert.Equal(t, id, inventory.ID)
	assert.NotNil(t, inventory.Product)
	assert.NotNil(t, inventory.Warehouse)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepository_UpdateStock(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()
	version := 1
	delta := 5

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "inventory" SET "quantity"=quantity + $1,"version"=version + 1 WHERE id = $2 AND version = $3`)).
		WithArgs(delta, id, version).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateStock(context.Background(), id, delta, version)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInventoryRepository_UpdateReservedStock(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewInventoryRepository(gormDB)

	id := uuid.New()
	version := 1
	delta := 2

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "inventory" SET "reserved_quantity"=reserved_quantity + $1,"version"=version + 1 WHERE id = $2 AND version = $3`)).
		WithArgs(delta, id, version).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateReservedStock(context.Background(), id, delta, version)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
