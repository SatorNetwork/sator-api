package google_firebase

import (
	"context"

	"firebase.google.com/go/messaging"
)

//go:generate mockgen -destination=mock_client.go -package=google_firebase github.com/SatorNetwork/sator-api/lib/google_firebase AppInterface
type AppInterface interface {
	Messaging(ctx context.Context) (*messaging.Client, error)
}
