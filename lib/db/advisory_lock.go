package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
)

type advisoryLocks struct {
	db *sql.DB
}

func NewAdvisoryLocks(db *sql.DB) GetLocker {
	return &advisoryLocks{
		db: db,
	}
}

func (a *advisoryLocks) GetLock(ctx context.Context, id string) (Locker, error) {
	idInt64 := convertStringToInt64(id)
	lock, err := NewLock(ctx, idInt64, a.db)
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

// Lock implements the Locker interface.
type Lock struct {
	id   int64
	conn *sql.Conn
}

// Lock obtains exclusive session level advisory lock if available.
// It’s similar to WaitAndLock, except it will not wait for the lock to become available.
// It will either obtain the lock and return true, or return false if the lock cannot be acquired immediately.
func (l *Lock) Lock(ctx context.Context) (bool, error) {
	result := false
	sqlQuery := "SELECT pg_try_advisory_lock($1)"
	err := l.conn.QueryRowContext(ctx, sqlQuery, l.id).Scan(&result)
	return result, err
}

// WaitAndLock obtains exclusive session level advisory lock.
// If another session already holds a lock on the same resource identifier, this function will wait until the resource becomes available.
// Multiple lock requests stack, so that if the resource is locked three times it must then be unlocked three times.
func (l *Lock) WaitAndLock(ctx context.Context) error {
	sqlQuery := "SELECT pg_advisory_lock($1)"
	_, err := l.conn.ExecContext(ctx, sqlQuery, l.id)
	return err
}

// Unlock releases the lock and DB connection.
func (l *Lock) Unlock(ctx context.Context) error {
	sqlQuery := "SELECT pg_advisory_unlock($1)"
	_, err := l.conn.ExecContext(ctx, sqlQuery, l.id)
	if err != nil {
		return err
	}
	// Returns the connection to the connection pool
	return l.conn.Close()
}

// NewLock returns a Lock with *sql.Conn
func NewLock(ctx context.Context, id int64, db *sql.DB) (Lock, error) {
	// Obtain a connection from the DB connection pool and store it and use it for lock and unlock operations
	conn, err := db.Conn(ctx)
	if err != nil {
		return Lock{}, err
	}
	return Lock{id: id, conn: conn}, nil
}
