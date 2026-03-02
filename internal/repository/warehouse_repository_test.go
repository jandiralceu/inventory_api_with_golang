package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()
	warehouse := &models.Warehouse{
		ID:          id,
		Name:        "Main Warehouse",
		Slug:        "main-warehouse",
		Code:        "WH01",
		Description: "Primary storage",
		IsActive:    true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "warehouses"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), warehouse)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Update_Success(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()
	warehouse := &models.Warehouse{
		ID:   id,
		Name: "Updated Name",
		Slug: "updated-slug",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "warehouses" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), warehouse)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Update_NotFound(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()
	warehouse := &models.Warehouse{
		ID:   id,
		Name: "Some Name",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "warehouses" SET`)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), warehouse)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Delete_Success(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "warehouses" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_Delete_NotFound(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "warehouses" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWarehouseRepository_FindByID(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "slug", "code"}).
		AddRow(id, "Main Warehouse", "main-warehouse", "WH01")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE id = $1`)).
		WithArgs(id, 1).
		WillReturnRows(rows)

	res, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, id, res.ID)
		assert.Equal(t, "WH01", res.Code)
	}
}

func TestWarehouseRepository_FindByCode(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	code := "WH01"
	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "code"}).
		AddRow(id, "Main Warehouse", code)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE code = $1`)).
		WithArgs(code, 1).
		WillReturnRows(rows)

	res, err := repo.FindByCode(context.Background(), code)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, code, res.Code)
	}
}

func TestWarehouseRepository_FindAll(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewWarehouseRepository(gormDB)

	filter := WarehouseListFilter{
		Name: "Main",
		Pagination: PaginationParams{
			Page:  1,
			Limit: 10,
		},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "warehouses" WHERE name ILIKE $1`)).
		WithArgs("%Main%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "warehouses" WHERE name ILIKE $1 ORDER BY created_at DESC LIMIT $2`)).
		WithArgs("%Main%", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(uuid.New(), "Main Warehouse"))

	results, total, err := repo.FindAll(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}
