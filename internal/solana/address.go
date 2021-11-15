package solana

import (
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
)

func IsSolanaAddress(addr string) bool {
	pk, err := base58.Decode(addr)
	if err != nil {
		return false
	}

	if len(pk) != common.PublicKeyLength {
		return false
	}

	return true
}