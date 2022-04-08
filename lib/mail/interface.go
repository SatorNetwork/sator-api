package mail

import (
	"context"
)

//go:generate mockgen -destination=mock_client.go -package=mail github.com/SatorNetwork/sator-api/lib/mail Interface
type (
	Interface interface {
		SendVerificationCode(_ context.Context, email, otp string) error
		SendResetPasswordCode(_ context.Context, email, otp string) error
		SendDestroyAccountCode(_ context.Context, email, otp string) error
		SendInvitation(_ context.Context, email, invitedBy string) error
	}
)
