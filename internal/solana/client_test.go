package solana_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/SatorNetwork/sator-api/internal/solana"

	"github.com/portto/solana-go-sdk/assotokenprog"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
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

/*func TestTest (t *testing.T) {
	ctx := context.Background()
	c := client.NewClient(client.DevnetRPCEndpoint)

	privateKey := []byte{210,26,212,148,51,216,254,151,70,177,14,51,24,82,207,128,222,200,188,175,33,76,112,231,169,182,77,195,227,87,28,143,188,216,244,205,52,123,229,204,198,210,96,123,243,180,194,108,175,124,122,229,104,69,144,208,62,200,100,237,132,82,251,23}
	mintAuth := "DiBXR49zpsJtgGd4hCZitmuWX19ko3Dc48yXBUgqKTxA"
	Address := "CizSaMmnZymceaDTPcNdXgKEpLarCQDvtAkAZA2tSE2u"

	feePayer := types.NewAccount()
	_, err := c.RequestAirdrop(ctx, feePayer.PublicKey.ToBase58(), 1000000000)
	require.NoError(t, err)

	testAccount := types.NewAccount()
	require.NoError(t, err)

	ata, _, err := common.FindAssociatedTokenAddress(testAccount.PublicKey, common.PublicKeyFromString(Address))
	require.NoError(t, err)

	res, err := c.GetRecentBlockhash(context.Background())
	require.NoError(t, err)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			tokenprog.TransferChecked(
				common.PublicKeyFromString(mintAuth),
				ata,
				common.PublicKeyFromString(Address),
				common.PublicKeyFromString(mintAuth),
				[]common.PublicKey{},
				10000,
				9,
			),
		},
		Signers:         []types.Account{feePayer, types.AccountFromPrivateKeyBytes(privateKey)},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err := c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	b, err := c.GetBalance(ctx, testAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, b)
}*/

func TestTest(t *testing.T) {
	ctx := context.Background()
	c := client.NewClient(client.DevnetRPCEndpoint)

	// Create mint
	feePayer := types.Account{
		PublicKey:  common.PublicKeyFromString("DiBXR49zpsJtgGd4hCZitmuWX19ko3Dc48yXBUgqKTxA"),
		PrivateKey: ed25519.PrivateKey{210, 26, 212, 148, 51, 216, 254, 151, 70, 177, 14, 51, 24, 82, 207, 128, 222, 200, 188, 175, 33, 76, 112, 231, 169, 182, 77, 195, 227, 87, 28, 143, 188, 216, 244, 205, 52, 123, 229, 204, 198, 210, 96, 123, 243, 180, 194, 108, 175, 124, 122, 229, 104, 69, 144, 208, 62, 200, 100, 237, 132, 82, 251, 23},
	}

	alice := types.NewAccount()
	mint := types.NewAccount()

	rentExemptionBalance, err := c.GetMinimumBalanceForRentExemption(ctx, tokenprog.MintAccountSize)
	require.NoError(t, err)

	res, err := c.GetRecentBlockhash(ctx)
	require.NoError(t, err)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				mint.PublicKey,
				common.TokenProgramID,
				rentExemptionBalance,
				tokenprog.MintAccountSize,
			),
			tokenprog.InitializeMint(
				8,
				mint.PublicKey,
				alice.PublicKey,
				common.PublicKey{},
			),
		},
		Signers:         []types.Account{feePayer, mint},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err := c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	// Random Token Account

	aliceTokenAccount := types.NewAccount()
	fmt.Println("alice token account:", aliceTokenAccount.PublicKey.ToBase58())

	rentExemptionBalance, err = c.GetMinimumBalanceForRentExemption(ctx, tokenprog.TokenAccountSize)
	require.NoError(t, err)

	res, err = c.GetRecentBlockhash(ctx)
	require.NoError(t, err)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				feePayer.PublicKey,
				aliceTokenAccount.PublicKey,
				common.TokenProgramID,
				rentExemptionBalance,
				tokenprog.TokenAccountSize,
			),
			tokenprog.InitializeAccount(
				aliceTokenAccount.PublicKey,
				mint.PublicKey,
				alice.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer, aliceTokenAccount},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err = c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	// Associated Token Account

	ata, _, err := common.FindAssociatedTokenAddress(alice.PublicKey, mint.PublicKey)
	require.NoError(t, err)
	fmt.Println("ata:", ata.ToBase58())

	res, err = c.GetRecentBlockhash(ctx)
	require.NoError(t, err)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			assotokenprog.CreateAssociatedTokenAccount(
				feePayer.PublicKey,
				alice.PublicKey,
				mint.PublicKey,
			),
		},
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err = c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	// Mint To

	res, err = c.GetRecentBlockhash(ctx)
	require.NoError(t, err)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			tokenprog.MintToChecked(
				mint.PublicKey,
				aliceTokenAccount.PublicKey,
				alice.PublicKey,
				[]common.PublicKey{},
				1e8,
				8,
			),
		},
		Signers:         []types.Account{feePayer, alice},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err = c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	// Transfer

	res, err = c.GetRecentBlockhash(ctx)
	require.NoError(t, err)

	rawTx, err = types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			tokenprog.TransferChecked(
				aliceTokenAccount.PublicKey,
				ata,
				mint.PublicKey,
				alice.PublicKey,
				[]common.PublicKey{},
				1e8,
				8,
			),
		},
		Signers:         []types.Account{feePayer, alice},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	require.NoError(t, err)

	txhash, err = c.SendRawTransaction(ctx, rawTx)
	require.NoError(t, err)
	require.NotNil(t, txhash)

	// Check balance?

	balance, err := c.GetBalance(ctx, alice.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance)
	balance, err = c.GetBalance(ctx, aliceTokenAccount.PublicKey.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance)
	balance, err = c.GetBalance(ctx, ata.ToBase58())
	require.NoError(t, err)
	require.NotNil(t, balance)
}
