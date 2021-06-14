package solana

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

type (
	Client struct {
		solana   *client.Client
		decimals uint8
		mltpl    uint64
	}
)

// New creates new solana client wrapper
func New(endpoint string) *Client {
	return &Client{
		solana:   client.NewClient(endpoint),
		decimals: 9,
		mltpl:    1e9,
	}
}

// NewAccount generates account keypair
func (c *Client) NewAccount() types.Account {
	return types.NewAccount()
}

func (c *Client) PublicKeyFromString(pk string) common.PublicKey {
	return common.PublicKeyFromString(pk)
}

func (c *Client) AccountFromPrivatekey(pk []byte) types.Account {
	return types.AccountFromPrivateKeyBytes(pk)
}

// RequestAirdrop working only in test and dev environment
func (c *Client) RequestAirdrop(pubKey string, amount float64) (string, error) {
	if amount > 10 {
		log.Printf("requested airdrop is too large %f, max: 10 SOL", amount)
		amount = 10
	}
	txhash, err := c.solana.RequestAirdrop(
		context.Background(),
		pubKey,
		uint64(amount*float64(c.mltpl)),
	)
	if err != nil {
		return "", fmt.Errorf("could not request airdrop: %w", err)
	}
	return txhash, nil
}

// SendTransaction sends transaction ans returns transaction hash
func (c *Client) SendTransaction(feePayer, signer types.Account, instructions ...types.Instruction) (string, error) {
	res, err := c.solana.GetRecentBlockhash(context.Background())
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions:    instructions,
		Signers:         []types.Account{feePayer, signer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		return "", fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, nil
}

// GetAccountBalanceSOL returns account's SOL balance
func (c *Client) GetAccountBalanceSOL(accPubKey string) (float64, error) {
	balance, err := c.solana.GetBalance(context.Background(), accPubKey)
	if err != nil {
		return 0, fmt.Errorf("could not get account balance: %w", err)
	}

	return float64(balance) / float64(1e9), nil
}

// GetTokenAccountBalance returns token account's balance
func (c *Client) GetTokenAccountBalance(accPubKey string) (float64, error) {
	accBalance, err := c.solana.GetTokenAccountBalance(context.Background(), accPubKey, client.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("could not get token account balance: %w", err)
	}

	if accBalance.Amount == "" {
		return 0, nil
	}

	balance, err := strconv.ParseFloat(accBalance.UIAmountString, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse token account balance: %w", err)
	}

	return balance, nil
}
