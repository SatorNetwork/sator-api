//go:build !mock_solana

package client

import (
	"context"
	"fmt"
	"github.com/portto/solana-go-sdk/rpc"
	"log"
	"strings"

	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
)

type (
	Client struct {
		solana   *client.Client
		endpoint string
		decimals uint8
		mltpl    uint64
		config   Config
	}
)

// New creates new solana client wrapper
func New(endpoint string, config Config) lib_solana.Interface {
	return &Client{
		solana:   client.NewClient(endpoint),
		endpoint: endpoint,
		decimals: 9,
		mltpl:    1e9,
		config:   config,
	}
}

// NewAccount generates account keypair
func (c *Client) NewAccount() types.Account {
	return types.NewAccount()
}

func (c *Client) PublicKeyFromString(pk string) common.PublicKey {
	return common.PublicKeyFromString(pk)
}

func (c *Client) AccountFromPrivateKeyBytes(pk []byte) (types.Account, error) {
	return types.AccountFromBytes(pk)
}

func (c *Client) CheckPrivateKey(addr string, pk []byte) error {
	account, err := c.AccountFromPrivateKeyBytes(pk)
	if err != nil {
		return err
	}

	addrFromPk := account.PublicKey.ToBase58()
	if !strings.EqualFold(addrFromPk, addr) {
		return fmt.Errorf("CheckPrivateKey: want = %s, got = %s", addr, addrFromPk)
	}

	return nil
}

func (c *Client) deriveATAPublicKey(ctx context.Context, recipientPK, assetPK common.PublicKey) (common.PublicKey, error) {
	// Check if the given account is already ATA or not
	recipientAddr := recipientPK.ToBase58()
	resp, err := c.solana.GetAccountInfo(ctx, recipientAddr)
	if err != nil {
		return common.PublicKey{}, err
	}
	if resp.Owner == common.TokenProgramID.ToBase58() {
		// given recipient public key is already an SPL token account
		return recipientPK, nil
	}

	// Getting of the recipient ATA
	recipientAta, _, err := common.FindAssociatedTokenAddress(recipientPK, assetPK)
	if err != nil {
		return common.PublicKey{}, err
	}
	// Check if the ATA already created
	ataInfo, err := c.solana.GetAccountInfo(ctx, recipientAta.ToBase58())
	if err != nil {
		return common.PublicKey{}, err
	}
	if ataInfo.Owner == common.TokenProgramID.ToBase58() {
		// given recipient public key is already an SPL token account
		return recipientAta, nil
	}

	return common.PublicKey{}, ErrATANotCreated
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
func (c *Client) SendTransaction(ctx context.Context, feePayer, signer types.Account, instructions ...types.Instruction) (string, error) {
	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	rawTx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			Instructions:    instructions,
			RecentBlockhash: res.Blockhash,
		}),
		Signers: []types.Account{feePayer, signer},
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendTransaction(ctx, rawTx)
	if err != nil {
		return "", fmt.Errorf("could not send transaction: %w", err)
	}

	return txhash, nil
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
	accBalance, _, err := c.solana.GetTokenAccountBalanceWithConfig(ctx, accPubKey, rpc.GetTokenAccountBalanceConfig{
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return 0, fmt.Errorf("could not get token account balance: %w", err)
	}

	return float64(accBalance) / float64(c.mltpl), nil
}

func (c *Client) GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error) {
	accountPublicKey := common.PublicKeyFromString(accountAddr)
	assetPublicKey := common.PublicKeyFromString(assetAddr)
	accountAta, _, err := common.FindAssociatedTokenAddress(accountPublicKey, assetPublicKey)
	if err != nil {
		return 0, err
	}

	return c.GetTokenAccountBalance(ctx, accountAta.ToBase58())
}

// GetTransactions ...
func (c *Client) GetTransactions(ctx context.Context, accPubKey string) (txList []lib_solana.ConfirmedTransactionResponse, err error) {
	signatures, err := c.solana.GetSignaturesForAddressWithConfig(ctx, accPubKey, rpc.GetSignaturesForAddressConfig{
		Limit:      30,
		Commitment: rpc.CommitmentFinalized,
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

func (c *Client) GetTransactionsWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (txList []lib_solana.ConfirmedTransactionResponse, err error) {
	accountPublicKey := common.PublicKeyFromString(accountAddr)
	assetPublicKey := common.PublicKeyFromString(assetAddr)
	accountAta, _, err := common.FindAssociatedTokenAddress(accountPublicKey, assetPublicKey)
	if err != nil {
		return nil, err
	}

	return c.GetTransactions(ctx, accountAta.ToBase58())
}

func (c *Client) FindAssociatedTokenAddress(walletAddress, tokenMintAddress common.PublicKey) (common.PublicKey, int, error) {
	return common.FindAssociatedTokenAddress(walletAddress, tokenMintAddress)
}
