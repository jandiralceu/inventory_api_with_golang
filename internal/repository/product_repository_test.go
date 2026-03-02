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

func TestProductRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	id := uuid.New()
	categoryID := uuid.New()
	supplierID := uuid.New()
	product := &models.Product{
		ID:          id,
		Name:        "Test Product",
		Slug:        "test-product",
		SKU:         "SKU-123",
		Description: "A test product",
		Price:       99.99,
		CategoryID:  &categoryID,
		SupplierID:  &supplierID,
		IsActive:    true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products"`)).
		WithArgs(product.SKU, product.Slug, product.Name, product.Description, product.Price, product.CostPrice, product.CategoryID, product.SupplierID, product.IsActive, product.ReorderLevel, product.ReorderQuantity, product.WeightKg, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(product.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), product)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_FindByID_NotFound(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."id" = $1`)).
		WithArgs(id, 1).
		WillReturnError(context.DeadlineExceeded)

	product, err := repo.FindByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_FindByID_WithPreload(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	id := uuid.New()
	catID := uuid.New()
	supID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "slug", "sku", "category_id", "supplier_id"}).
		AddRow(id, "Test Product", "test-product", "SKU-123", catID, supID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."id" = $1`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
		WithArgs(catID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(catID, "Electronics"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "suppliers" WHERE "suppliers"."id" = $1`)).
		WithArgs(supID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(supID, "Global Corp"))

	product, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, id, product.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_FindBySKU(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	sku := "SKU-123"
	id := uuid.New()
	catID := uuid.New()
	supID := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "sku", "category_id", "supplier_id"}).
		AddRow(id, "Test Product", sku, catID, supID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE sku = $1`)).
		WithArgs(sku, 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."id" = $1`)).
		WithArgs(catID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(catID, "Electronics"))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "suppliers" WHERE "suppliers"."id" = $1`)).
		WithArgs(supID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(supID, "Global Corp"))

	product, err := repo.FindBySKU(context.Background(), sku)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, sku, product.SKU)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Update(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	id := uuid.New()
	product := &models.Product{
		ID:        id,
		Name:      "Updated Product",
		Slug:      "updated-product",
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), product)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepository_Delete(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewProductRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "products" WHERE "products"."id" = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
