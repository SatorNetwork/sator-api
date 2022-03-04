package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

type (
	Client struct {
		initSPLBalance              float64
		initSolBalance              float64
		errGiveAssetsWithAutoDerive error
		errGetTokenAccountBalance   error
	}
)

func New(initSPLBalance, initSolBalance float64,
	errGiveAssetsWithAutoDerive, errGetTokenAccountBalance error) *Client {

	return &Client{
		initSPLBalance:              initSPLBalance,
		initSolBalance:              initSolBalance,
		errGiveAssetsWithAutoDerive: errGiveAssetsWithAutoDerive,
		errGetTokenAccountBalance:   errGetTokenAccountBalance,
	}
}

func (c *Client) GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error) {
	if c.errGiveAssetsWithAutoDerive != nil {
		return "", c.errGiveAssetsWithAutoDerive
	}

	return fmt.Sprintf("mock_solana_tr_address_%d", time.Now().UnixNano()), nil
}

func (c *Client) GetTokenAccountBalance(ctx context.Context, accPubKey string) (float64, error) {
	if c.errGetTokenAccountBalance != nil {
		return 0, c.errGetTokenAccountBalance
	}

	return 100, nil
}

func (c *Client) RequestAirdrop(ctx context.Context, pubKey string, amount float64) (string, error) {
	return fmt.Sprintf("mock_solana_tr_address_%d", time.Now().UnixNano()), nil
}

func (c *Client) GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error) {
	return c.initSolBalance, nil
}

func (c *Client) CreateAccountWithATA(ctx context.Context, assetAddr, initAccAddr string, feePayer types.Account) (string, error) {
	return fmt.Sprintf("mock_solana_tr_address_%d", time.Now().UnixNano()), nil
}

func (c *Client) AccountFromPrivateKeyBytes(pk []byte) types.Account {
	return types.AccountFromPrivateKeyBytes(pk)
}

func (c *Client) CreateAsset(ctx context.Context, feePayer, issuer, asset types.Account) (string, error) {
	return fmt.Sprintf("mock_solana_tr_address_%d", time.Now().UnixNano()), nil
}

func (c *Client) IssueAsset(ctx context.Context, feePayer, issuer, asset types.Account, dest common.PublicKey, amount float64) (string, error) {
	return fmt.Sprintf("mock_solana_tr_address_%d", time.Now().UnixNano()), nil
}
