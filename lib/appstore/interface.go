package appstore

import (
	"context"

	"github.com/awa/go-iap/appstore"
)

//go:generate mockgen -destination=mock_client.go -package=appstore github.com/SatorNetwork/sator-api/lib/appstore Interface
type (
	Interface interface {
		Verify(ctx context.Context, reqBody appstore.IAPRequest, result interface{}) error
	}
)
