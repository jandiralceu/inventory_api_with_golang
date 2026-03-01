package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSupplierRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	id := uuid.New()
	supplier := &models.Supplier{
		ID:    id,
		Name:  "Supplier One",
		Slug:  "supplier-one",
		TaxID: "123456789",
		Address: models.Address{
			Street: "Main St",
			City:   "Tech City",
		},
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "suppliers"`)).
		WithArgs(supplier.Name, supplier.Slug, supplier.Description, supplier.TaxID, supplier.Email, supplier.Phone, sqlmock.AnyArg(), supplier.ContactPerson, supplier.IsActive, supplier.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(supplier.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), supplier)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierRepository_FindByID(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "slug", "tax_id"}).
		AddRow(id, "Supplier One", "supplier-one", "123456789")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "suppliers" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	supplier, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, supplier)
	assert.Equal(t, id, supplier.ID)
	assert.Equal(t, "Supplier One", supplier.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierRepository_FindByID_NotFound(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "suppliers" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	supplier, err := repo.FindByID(context.Background(), id)
	assert.Error(t, err)
	assert.Nil(t, supplier)
	// GORM's First() returns gorm.ErrRecordNotFound when no rows are found.
	// Our repository maps this to apperrors.ErrNotFound.
	// However, depending on how GORM handles sql.ErrNoRows from mock, it might return it directly or wrap it.
	// To be safe and test our mapper, we check for apperrors.ErrNotFound.
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestSupplierRepository_Update(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	id := uuid.New()
	supplier := &models.Supplier{
		ID:        id,
		Name:      "Supplier Updated",
		Slug:      "supplier-updated",
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "suppliers" SET`)).
		WithArgs(supplier.Name, supplier.Slug, supplier.Description, supplier.TaxID, supplier.Email, supplier.Phone, sqlmock.AnyArg(), supplier.ContactPerson, supplier.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg(), supplier.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), supplier)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierRepository_Delete(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "suppliers" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSupplierRepository_FindAll(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewSupplierRepository(gormDB)

	filter := SupplierListFilter{
		Name: "Supplier",
		Pagination: PaginationParams{
			Page:  1,
			Limit: 10,
			Sort:  "name",
			Order: "asc",
		},
	}

	// Mock Count
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "suppliers" WHERE name ILIKE $1`)).
		WithArgs("%Supplier%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock Data
	rows := sqlmock.NewRows([]string{"id", "name", "slug"}).
		AddRow(uuid.New(), "Supplier One", "supplier-one")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "suppliers" WHERE name ILIKE $1 ORDER BY name asc LIMIT $2`)).
		WithArgs("%Supplier%", 10).
		WillReturnRows(rows)

	suppliers, total, err := repo.FindAll(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, suppliers, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
