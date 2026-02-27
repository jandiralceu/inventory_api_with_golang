package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"gorm.io/gorm"
)

// mapDatabaseError translates database-specific errors (GORM, Postgres)
// into standardized application errors defined in the apperrors package.
func mapDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Postgres error code 23505 represents a unique constraint violation.
		switch pgErr.Code {
		case "23505": // unique_violation
			return apperrors.ErrConflict
		}
	}

	return err
}
