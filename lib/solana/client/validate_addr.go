package client

import (
	"fmt"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
)

// ValidateSolanaWalletAddr validates a Solana wallet address
func ValidateSolanaWalletAddr(addr string) error {
	d, err := base58.Decode(addr)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPublicKey, err)
	}

	if len(d) != common.PublicKeyLength {
		return fmt.Errorf("%w: wrong length", ErrInvalidPublicKey)
	}

	if _, err := new(edwards25519.Point).SetBytes(d); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidPublicKey, err)
	}

	return nil
}
