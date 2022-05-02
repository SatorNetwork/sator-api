//go:build mock_solana

package client

import (
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func New(endpoint string, config Config, exchangeRatesClient *exchange_rates_client.Client) lib_solana.Interface {
	m := mock.GetMockObject(mock.SolanaProvider)
	if m == nil {
		return nil
	}
	return m.(lib_solana.Interface)
}
