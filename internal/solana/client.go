package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

// Client holds all required fields to connect Solana network.
type Client struct {
	c client.Client
}

// New constructor for Solana client.
func New() *Client {
	return &Client{c: *client.NewClient(client.DevnetRPCEndpoint)}
}

// CreateAccount creates new account and requests airdrop to activate account in solana network, returns base58 public key and private key in []byte.
func (cl *Client) CreateAccount(ctx context.Context) (string, []byte, error) {
	account := types.NewAccount()

	_, err := cl.c.RequestAirdrop(context.Background(), account.PublicKey.ToBase58(), 1000000000)
	if err != nil {
		return "", nil, err
	}

	return account.PublicKey.ToBase58(), account.PrivateKey, nil
}

// GetBalance returns account's balance by base58 public key.
func (cl *Client) GetBalance(ctx context.Context, base58key string) (uint64, error) {
	return cl.c.GetBalance(ctx, base58key)
}

// SendTo creates transaction from root account to receiver amount in SOL (1000000000 =  1 SOL).
func (cl *Client) SendTo(ctx context.Context, receiverBase58Key string, amount uint64) error {
	res, err := cl.c.GetRecentBlockhash(context.Background())
	if err != nil {
		return fmt.Errorf("could not get recent block hash: %w", err)
	}

	// TODO: figure out what account to use as feePayer: atm using hardcoded from solana SDK guide.
	feePayerPK := []byte{57, 17, 193, 142, 252, 221, 81, 90, 60, 28, 93, 237, 212, 51, 95, 95, 41, 104, 221, 59, 13, 244, 54, 1, 79, 180, 120, 178, 81, 45, 46, 193, 142, 11, 237, 209, 82, 24, 36, 72, 7, 76, 66, 215, 44, 116, 17, 132, 252, 205, 47, 74, 57, 230, 36, 98, 119, 86, 11, 40, 71, 195, 47, 254}
	payingAccPK := []byte{210, 26, 212, 148, 51, 216, 254, 151, 70, 177, 14, 51, 24, 82, 207, 128, 222, 200, 188, 175, 33, 76, 112, 231, 169, 182, 77, 195, 227, 87, 28, 143, 188, 216, 244, 205, 52, 123, 229, 204, 198, 210, 96, 123, 243, 180, 194, 108, 175, 124, 122, 229, 104, 69, 144, 208, 62, 200, 100, 237, 132, 82, 251, 23}
	feePayer := types.AccountFromPrivateKeyBytes(feePayerPK)
	accountA := types.AccountFromPrivateKeyBytes(payingAccPK)

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.Transfer(
				accountA.PublicKey, // from
				common.PublicKeyFromString(receiverBase58Key), // to
				amount,
			),
		},
		Signers:         []types.Account{feePayer, accountA},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return fmt.Errorf("could not create raw transaction: %w", err)
	}

	_, err = cl.c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		return fmt.Errorf("could not send transaction: %w", err)
	}

	return nil
}
