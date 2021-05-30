package solana_test

import (
	"context"
	"reflect"
	"testing"

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

		// 9994429840
		balance, err = c.GetBalance(ctx, "DiBXR49zpsJtgGd4hCZitmuWX19ko3Dc48yXBUgqKTxA")
		require.NoError(t, err)
		require.NotNil(t, balance)

		// 9954
		balance, err = c.GetBalance(ctx, account.PublicKey.ToBase58())
		require.NoError(t, err)
		require.NotNil(t, balance)
	})
}

func TestClient_CreateAccount(t *testing.T) {
	type fields struct {
		c *client.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := &Client{
				c: tt.fields.c,
			}
			got, got1, err := cl.CreateAccount(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.CreateAccount() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Client.CreateAccount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
