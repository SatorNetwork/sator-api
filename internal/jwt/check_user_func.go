package jwt

import (
	"context"

	"github.com/google/uuid"
)

func CheckUser(isUserDisabled func(context.Context, uuid.UUID) (bool, error)) func(context.Context) error {
	return func(ctx context.Context) error {
		uid, err := UserIDFromContext(ctx)
		if err != nil {
			return ErrMissedUserID
		}

		if yes, err := isUserDisabled(ctx, uid); err != nil || yes {
			return ErrUserIsDisabled
		}

		return nil
	}
}
