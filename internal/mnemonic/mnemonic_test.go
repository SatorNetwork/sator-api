package mnemonic_test

import (
	"testing"

	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/assert"

	mnemonic2 "github.com/SatorNetwork/sator-api/internal/mnemonic"
)

func TestMnemonic(t *testing.T) {
	initAcc := types.NewAccount()

	mnemonic, err := mnemonic2.GenerateMnemonicPassphrase(initAcc.PrivateKey)
	assert.NoError(t, err)
	assert.NotEqual(t, mnemonic, "")

	newAccount, err := mnemonic2.KeyFromMnemonic(mnemonic)
	assert.NoError(t, err)
	assert.Equal(t, initAcc, newAccount)
}
