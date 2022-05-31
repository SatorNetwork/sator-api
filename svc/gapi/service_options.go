package gapi

import (
	"database/sql"
)

// WithDB sets the database connection
func WithDB(db *sql.DB) ServiceOption {
	return func(s *Service) {
		s.db = db
	}
}
