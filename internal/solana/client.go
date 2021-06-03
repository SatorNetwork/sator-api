package solana

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/types"
)

// Client holds all required fields to connect Solana network.
type Client struct {
	c             client.Client
	feePayerPK    []byte
	rewardPayerPK []byte
}

// New constructor for Solana client.
func New(endpoint string, feePayerPK, rewardPayerPK []byte) *Client {
	return &Client{c: *client.NewClient(endpoint), feePayerPK: feePayerPK, rewardPayerPK: rewardPayerPK}
}

// CreateAccount creates new account and requests airdrop to activate account in solana network, returns base58 public key and private key in []byte.
func (cl *Client) CreateAccount(ctx context.Context) (string, []byte, error) {
	account := types.NewAccount()

	// workaround to avoid internal error while user signing up
	for i := 0; i < 3; i++ {
		_, err := cl.c.RequestAirdrop(ctx, account.PublicKey.ToBase58(), 1000000000)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(100*(i+1)))
		if i == 2 {
			log.Printf("could not request airdrop for account %s", account.PublicKey.ToBase58())
		}
	}

	return account.PublicKey.ToBase58(), account.PrivateKey, nil
}

// GetBalance returns account's balance by base58 public key.
func (cl *Client) GetBalance(ctx context.Context, base58key string) (uint64, error) {
	return cl.c.GetBalance(ctx, base58key)
}

// SendTo creates transaction from root account to receiver amount in SOL (1000000000 =  1 SOL), returns transaction hash.
func (cl *Client) SendTo(ctx context.Context, receiverBase58Key string, amount uint64) (string, error) {
	res, err := cl.c.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	feePayer := types.AccountFromPrivateKeyBytes(cl.feePayerPK)
	accountA := types.AccountFromPrivateKeyBytes(cl.rewardPayerPK)

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
		return "", fmt.Errorf("could not create raw transaction: %w", err)
	}

	txHash, err := cl.c.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", fmt.Errorf("could not send transaction: %w", err)
	}

	return txHash, nil
}

// GetTransaction returns solana transaction by hash.
func (cl *Client) GetTransaction(ctx context.Context, txHash string) (string, error) {
	// TODO: figure out what to return from tx (can't find direct link atm)
	tx, err := cl.c.GetConfirmedTransaction(ctx, txHash)
	if err != nil {
		return "", fmt.Errorf("could not get transaction: %w", err)
	}
	log.Println(tx)

	return tx.Transaction.Message.RecentBlockhash, nil
}
