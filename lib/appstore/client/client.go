//go:build !mock_appstore

package client

import (
	"github.com/awa/go-iap/appstore"

	lib_appstore "github.com/SatorNetwork/sator-api/lib/appstore"
)

func New() lib_appstore.Interface {
	return appstore.New()
}
