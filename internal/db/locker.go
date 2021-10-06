package db

import (
	"context"

	"crypto/sha256"
	"database/sql"
	"encoding/binary"

	"github.com/allisson/go-pglock/v2"
)

type (
	GetLocker interface {
		GetLock(ctx context.Context, id string) (pglock.Locker, error)
	}
)

type advisoryLocks struct {
	db *sql.DB
}

func NewAdvisoryLocks(db *sql.DB) GetLocker {
	return &advisoryLocks{
		db: db,
	}
}

func (a *advisoryLocks) GetLock(ctx context.Context, id string) (pglock.Locker, error) {
	idInt64 := convertStringToInt64(id)
	lock, err := pglock.NewLock(ctx, idInt64, a.db)
	if err != nil {
		return nil, err
	}

	return &lock, nil
}

func convertStringToInt64(id string) int64 {
	idHash := sha256.Sum256([]byte(id))
	idUint64 := binary.BigEndian.Uint64(idHash[:8])
	idInt64 := int64(idUint64)

	return idInt64
}
