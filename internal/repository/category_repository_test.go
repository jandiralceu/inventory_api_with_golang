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

func TestCategoryRepository_Create(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewCategoryRepository(gormDB)

	id := uuid.New()
	category := &models.Category{
		ID:          id,
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic devices",
		IsActive:    true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories"`)).
		WithArgs(category.Name, category.Slug, category.Description, category.ParentID, category.IsActive, category.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(category.ID, time.Now(), time.Now()))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), category)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_FindByID(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewCategoryRepository(gormDB)

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "slug"}).
		AddRow(id, "Electronics", "electronics")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE id = $1`)).
		WithArgs(id, 1). // GORM adds LIMIT 1
		WillReturnRows(rows)

	category, err := repo.FindByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, id, category.ID)
	assert.Equal(t, "Electronics", category.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_FindBySlug(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewCategoryRepository(gormDB)

	slug := "electronics"
	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "name", "slug"}).
		AddRow(id, "Electronics", slug)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE slug = $1`)).
		WithArgs(slug, 1). // GORM adds LIMIT 1
		WillReturnRows(rows)

	category, err := repo.FindBySlug(context.Background(), slug)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, slug, category.Slug)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_Update(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewCategoryRepository(gormDB)

	id := uuid.New()
	category := &models.Category{
		ID:          id,
		Name:        "Electronics Updated",
		Slug:        "electronics-updated",
		Description: "Updated description",
		IsActive:    true,
		UpdatedAt:   time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET`)).
		WithArgs(category.Name, category.Slug, category.Description, category.ParentID, category.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg(), category.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), category)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_Delete(t *testing.T) {
	gormDB, mock, _ := setupTestDB(t)
	repo := NewCategoryRepository(gormDB)

	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), id)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
