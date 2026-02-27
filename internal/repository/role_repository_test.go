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

// =====================
// Create Tests
// =====================

func TestRoleRepositoryCreateSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	role := &models.Role{
		Name:        "Admin",
		Description: "Administrator role",
	}

	roleID := uuid.New()
	now := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "roles" ("name","description") VALUES ($1,$2) RETURNING "id","created_at"`)).
		WithArgs("admin", role.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(roleID, now))
	mock.ExpectCommit()

	result, err := repo.Create(context.Background(), role)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryCreateDuplicateError(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	role := &models.Role{
		Name:        "Admin",
		Description: "Administrator role",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "roles"`)).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	result, err := repo.Create(context.Background(), role)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FindAll Tests
// =====================

func TestRoleRepositoryFindAllSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID1 := uuid.New()
	roleID2 := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at"}).
		AddRow(roleID1, "Admin", "Administrator role", now).
		AddRow(roleID2, "Editor", "Editor role", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
		WillReturnRows(rows)

	result, err := repo.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Admin", result[0].Name)
	assert.Equal(t, "Editor", result[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryFindAllEmpty(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
		WillReturnRows(rows)

	result, err := repo.FindAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryFindAllError(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
		WillReturnError(gorm.ErrInvalidDB)

	result, err := repo.FindAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// FindByID Tests
// =====================

func TestRoleRepositoryFindByIDSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at"}).
		AddRow(roleID, "Admin", "Administrator role", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1`)).
		WithArgs(roleID, 1).
		WillReturnRows(rows)

	role, err := repo.FindByID(context.Background(), roleID)

	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, roleID, role.ID)
	assert.Equal(t, "Admin", role.Name)
	assert.Equal(t, "Administrator role", role.Description)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryFindByIDNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1`)).
		WithArgs(roleID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	role, err := repo.FindByID(context.Background(), roleID)

	assert.Error(t, err)
	assert.Nil(t, role)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// =====================
// Delete Tests
// =====================

func TestRoleRepositoryDeleteSuccess(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), roleID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryDeleteNotFound(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
		WithArgs(roleID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), roleID)

	assert.Error(t, err)
	assert.Equal(t, apperrors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepositoryDeleteError(t *testing.T) {
	gormDB, mock, db := setupTestDB(t)
	defer db.Close()

	repo := NewRoleRepository(gormDB)

	roleID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE id = $1`)).
		WithArgs(roleID).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	err := repo.Delete(context.Background(), roleID)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
