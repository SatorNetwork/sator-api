package db

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// IsNotFoundError determines if an error is sql.ErrNoRows
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// IsDuplicateError determines if an error is pq: duplicate key value violates unique constraint error
func IsDuplicateError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code.Name() == "unique_violation" {
			return true
		}
	}
	return false
}
