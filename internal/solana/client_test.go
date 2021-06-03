package solana_test

import (
	"context"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"testing"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	c := client.NewClient(client.DevnetRPCEndpoint)
	account := types.NewAccount()

	t.Run("test get balance", func(t *testing.T) {
		balance, err := c.GetBalance(ctx, account.PublicKey.ToBase58())
		require.NoError(t, err)
		require.NotNil(t, balance)
	})

	t.Run("test send transaction", func(t *testing.T) {
		res, err := c.GetRecentBlockhash(context.Background())
		require.NoError(t, err)

		feePayer := types.AccountFromPrivateKeyBytes([]byte{57, 17, 193, 142, 252, 221, 81, 90, 60, 28, 93, 237, 212, 51, 95, 95, 41, 104, 221, 59, 13, 244, 54, 1, 79, 180, 120, 178, 81, 45, 46, 193, 142, 11, 237, 209, 82, 24, 36, 72, 7, 76, 66, 215, 44, 116, 17, 132, 252, 205, 47, 74, 57, 230, 36, 98, 119, 86, 11, 40, 71, 195, 47, 254})
		accountA := types.AccountFromPrivateKeyBytes([]byte{210, 26, 212, 148, 51, 216, 254, 151, 70, 177, 14, 51, 24, 82, 207, 128, 222, 200, 188, 175, 33, 76, 112, 231, 169, 182, 77, 195, 227, 87, 28, 143, 188, 216, 244, 205, 52, 123, 229, 204, 198, 210, 96, 123, 243, 180, 194, 108, 175, 124, 122, 229, 104, 69, 144, 208, 62, 200, 100, 237, 132, 82, 251, 23})

		rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
			Instructions: []types.Instruction{
				sysprog.Transfer(
					accountA.PublicKey, // from
					account.PublicKey,  // to
					10000,              // 1 SOL
				),
			},
			Signers:         []types.Account{feePayer, accountA},
			FeePayer:        feePayer.PublicKey,
			RecentBlockHash: res.Blockhash,
		})
		require.NoError(t, err)

		txSig, err := c.SendRawTransaction(context.Background(), rawTx)
		require.NoError(t, err)
		require.NotNil(t, txSig)

		balance, err := c.GetBalance(ctx, account.PublicKey.ToBase58())
		require.NoError(t, err)
		require.NotNil(t, balance)

		balance, err = c.GetBalance(ctx, "DiBXR49zpsJtgGd4hCZitmuWX19ko3Dc48yXBUgqKTxA")
		require.NoError(t, err)
		require.NotNil(t, balance)

		balance, err = c.GetBalance(ctx, account.PublicKey.ToBase58())
		require.NoError(t, err)
		require.NotNil(t, balance)
	})
}

func TestClient_CreateAccount(t *testing.T) {
	feePayer := []byte{57, 17, 193, 142, 252, 221, 81, 90, 60, 28, 93, 237, 212, 51, 95, 95, 41, 104, 221, 59, 13, 244, 54, 1, 79, 180, 120, 178, 81, 45, 46, 193, 142, 11, 237, 209, 82, 24, 36, 72, 7, 76, 66, 215, 44, 116, 17, 132, 252, 205, 47, 74, 57, 230, 36, 98, 119, 86, 11, 40, 71, 195, 47, 254}
	accountA := []byte{210, 26, 212, 148, 51, 216, 254, 151, 70, 177, 14, 51, 24, 82, 207, 128, 222, 200, 188, 175, 33, 76, 112, 231, 169, 182, 77, 195, 227, 87, 28, 143, 188, 216, 244, 205, 52, 123, 229, 204, 198, 210, 96, 123, 243, 180, 194, 108, 175, 124, 122, 229, 104, 69, 144, 208, 62, 200, 100, 237, 132, 82, 251, 23}

	pk58, pk, err := solana.New(client.DevnetRPCEndpoint, feePayer, accountA).CreateAccount(context.Background())
	if err != nil {
		t.Fatalf("could not create solana account: %v", err)
	}
	if pk58 == "" {
		t.Fatal("public key: expected base58 string, got empty")
	}
	if len(pk) < 1 {
		t.Fatal("private key is empty")
	}
}

func TestClient_CreateTokenTransfer(t *testing.T) {
	ctx := context.Background()
	c := client.NewClient(client.DevnetRPCEndpoint)
	newAcc := types.NewAccount()
	newAcc2 := types.NewAccount()
	newAcc3 := types.NewAccount()

	//fromAcc := types.NewAccount()
	//ownerAcc := types.NewAccount()

	/*feePayer := []byte{57, 17, 193, 142, 252, 221, 81, 90, 60, 28, 93, 237, 212, 51, 95, 95, 41, 104, 221, 59, 13, 244, 54, 1, 79, 180, 120, 178, 81, 45, 46, 193, 142, 11, 237, 209, 82, 24, 36, 72, 7, 76, 66, 215, 44, 116, 17, 132, 252, 205, 47, 74, 57, 230, 36, 98, 119, 86, 11, 40, 71, 195, 47, 254}
	accountA := []byte{210, 26, 212, 148, 51, 216, 254, 151, 70, 177, 14, 51, 24, 82, 207, 128, 222, 200, 188, 175, 33, 76, 112, 231, 169, 182, 77, 195, 227, 87, 28, 143, 188, 216, 244, 205, 52, 123, 229, 204, 198, 210, 96, 123, 243, 180, 194, 108, 175, 124, 122, 229, 104, 69, 144, 208, 62, 200, 100, 237, 132, 82, 251, 23}
	key1 := common.PublicKeyFromString("DiBXR49zpsJtgGd4hCZitmuWX19ko3Dc48yXBUgqKTxA")*/

	_ = sysprog.CreateAccount(newAcc.PublicKey, newAcc2.PublicKey, newAcc3.PublicKey, 1000000, 1000000000)
	_ = tokenprog.InitializeMint(5, newAcc.PublicKey, newAcc2.PublicKey, newAcc3.PublicKey)
	_ = tokenprog.MintTo(newAcc.PublicKey, newAcc2.PublicKey, newAcc3.PublicKey, []common.PublicKey{newAcc.PublicKey, newAcc2.PublicKey}, 1000)
	_ = tokenprog.Transfer(newAcc.PublicKey, newAcc2.PublicKey, newAcc3.PublicKey, []common.PublicKey{newAcc.PublicKey, newAcc2.PublicKey}, 1000)

	balance, err := c.GetBalance(ctx, newAcc.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance)
	balance2, err := c.GetBalance(ctx, newAcc2.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance2)
	balance3, err := c.GetBalance(ctx, newAcc3.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance3)
}