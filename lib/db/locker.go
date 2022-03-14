package db

import (
	"context"
)

type (
	GetLocker interface {
		GetLock(ctx context.Context, id string) (Locker, error)
	}

	Locker interface {
		Lock(ctx context.Context) (bool, error)
		WaitAndLock(ctx context.Context) error
		Unlock(ctx context.Context) error
	}
)
