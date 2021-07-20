package solana

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

type (
	Client struct {
		solana   *client.Client
		endpoint string
		decimals uint8
		mltpl    uint64
		FeePayer types.Account
		Asset    types.Account
		Issuer   types.Account
	}
)

// New creates new solana client wrapper
func New(endpoint string, feePayer, asset, issuer types.Account) *Client {
	return &Client{
		solana:   client.NewClient(endpoint),
		endpoint: endpoint,
		decimals: 9,
		mltpl:    1e9,
		FeePayer: feePayer,
		Issuer:   issuer,
		Asset:    asset,
	}
}

// NewAccount generates account keypair
func (c *Client) NewAccount() types.Account {
	return types.NewAccount()
}

func (c *Client) PublicKeyFromString(pk string) common.PublicKey {
	return common.PublicKeyFromString(pk)
}

func (c *Client) AccountFromPrivateKey(pk []byte) types.Account {
	return types.AccountFromPrivateKeyBytes(pk)
}

// RequestAirdrop working only in test and dev environment
func (c *Client) RequestAirdrop(ctx context.Context, pubKey string, amount float64) (string, error) {
	if amount > 10 {
		log.Printf("requested airdrop is too large %f, max: 10 SOL", amount)
		amount = 10
	}
	txhash, err := c.solana.RequestAirdrop(
		ctx,
		pubKey,
		uint64(amount*float64(c.mltpl)),
	)
	if err != nil {
		return "", fmt.Errorf("could not request airdrop: %w", err)
	}
	return txhash, nil
}

// SendTransaction sends transaction ans returns transaction hash
func (c *Client) SendTransaction(ctx context.Context, signer types.Account, instructions ...types.Instruction) (txhash string, err error) {
	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions:    instructions,
		Signers:         []types.Account{c.FeePayer, signer},
		FeePayer:        c.FeePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	for i := 0; i < 3; i++ {
		txhash, err = c.solana.SendRawTransaction(ctx, rawTx)
		if err != nil {
			time.Sleep(10*time.Second)
			continue
		} else {
			break
		}
	}

	return txhash, err
}

// GetAccountBalanceSOL returns account's SOL balance
func (c *Client) GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error) {
	balance, err := c.solana.GetBalance(ctx, accPubKey)
	if err != nil {
		return 0, fmt.Errorf("could not get account balance: %w", err)
	}

	return float64(balance) / 1e9, nil
}

// GetTokenAccountBalance returns token account's balance
func (c *Client) GetTokenAccountBalance(ctx context.Context, accPubKey string) (float64, error) {
	accBalance, err := c.solana.GetTokenAccountBalance(ctx, accPubKey, client.CommitmentFinalized)
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

// GetTransactions ...
func (c *Client) GetTransactions(ctx context.Context, accPubKey string) (txList []ConfirmedTransactionResponse, err error) {
	signatures, err := c.solana.GetConfirmedSignaturesForAddress(ctx, accPubKey, client.GetConfirmedSignaturesForAddressConfig{
		Limit:      30,
		Commitment: client.CommitmentFinalized,
	})
	if err != nil {
		return nil, err
	}

	for _, signature := range signatures {
		tx, err := c.GetConfirmedTransactionForAccount(ctx, accPubKey, signature.Signature)
		if err != nil {
			return nil, err
		}

		txList = append(txList, tx)
	}

	return txList, nil
}
