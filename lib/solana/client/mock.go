//go:build mock_solana

package client

import (
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func New(endpoint string, config Config) lib_solana.Interface {
	m := mock.GetMockObject(mock.SolanaProvider)
	if m == nil {
		return nil
	}
	return m.(lib_solana.Interface)
}
