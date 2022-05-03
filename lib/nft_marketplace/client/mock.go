//go:build mock_nft_marketplace

package client

import (
	lib_nft_marketplace "github.com/SatorNetwork/sator-api/lib/nft_marketplace"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func New(serverHost string, serverPort int) lib_nft_marketplace.Interface {
	m := mock.GetMockObject(mock.NftMarketplaceProvider)
	if m == nil {
		return nil
	}
	return m.(lib_nft_marketplace.Interface)
}
