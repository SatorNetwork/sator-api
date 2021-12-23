package mnemonic

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSLIP10Compatibility(t *testing.T) {
	mnemonic := "diagram another jealous will cost ship goose blind elevator anxiety crazy cheese " +
		"cherry jeans rhythm february fat broom tattoo artwork cluster damp maple scorpion"
	account, err := AccountFromMnemonic(mnemonic)
	require.NoError(t, err)

	{
		expectedPrivateKeyHex := "623c0c7fbdd49b93a33aef2a1eada0f1f9ee7d06f958194ed8a7a1fa6b76d47f" +
			"7541f1271fecbb9fad2501077b20779d0fc5448c45fcd549ac7c2ba81cf676b0"
		actualPrivateKeyHex := hex.EncodeToString(account.PrivateKey)
		require.Equal(t, expectedPrivateKeyHex, actualPrivateKeyHex)
	}

	{
		expectedAddr := "8tj2AYrV3bNHaayZuTiQs5vShJH57PtnBsDYJT7QBEK9"
		actualAddr := account.PublicKey.ToBase58()
		require.Equal(t, expectedAddr, actualAddr)
	}
}

func TestMnemonic(t *testing.T) {
	mnemonic, err := NewMnemonic()
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)

	account, err := AccountFromMnemonic(mnemonic)
	require.NoError(t, err)
	require.NotNil(t, account)
	require.NotNil(t, account.PrivateKey)
	require.NotNil(t, account.PublicKey)
}
