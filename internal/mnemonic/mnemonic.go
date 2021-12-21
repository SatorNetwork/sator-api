package mnemonic

import (
	"github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ed25519"
)

func GenerateMnemonicPassphrase(privateKey ed25519.PrivateKey) (string, error) {
	return bip39.NewMnemonic(privateKey.Seed())
}

func KeyFromMnemonic(mnemonic string) (types.Account, error) {
	entropy, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return types.Account{}, err
	}

	key := ed25519.NewKeyFromSeed(entropy)
	return types.AccountFromPrivateKeyBytes(key), nil
}
