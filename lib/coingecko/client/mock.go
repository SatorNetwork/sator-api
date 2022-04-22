//go:build mock_coingecko

package client

import (
	lib_coingecko "github.com/SatorNetwork/sator-api/lib/coingecko"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func NewCoingeckoClient() lib_coingecko.Interface {
	m := mock.GetMockObject(mock.CoingeckoProvider)
	if m == nil {
		return nil
	}
	return m.(lib_coingecko.Interface)
}
