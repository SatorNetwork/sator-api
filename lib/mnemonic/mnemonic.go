package mnemonic

import (
	"github.com/anytypeio/go-slip10"
	"github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func AccountFromMnemonic(mnemonic string) (*types.Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	node, err := slip10.DeriveForPath("m/44'/501'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	_, privateKey := node.Keypair()
	account, err := types.AccountFromBytes(privateKey)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
