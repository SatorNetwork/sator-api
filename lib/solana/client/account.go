package client

import (
	"context"
	"github.com/portto/solana-go-sdk/rpc"
)

func (c *Client) GetTokenAccountBalanceCall(ctx context.Context, base58Addr string) (*rpc.GetAccountInfoResultValue, error) {
	var res rpc.GetAccountInfoResponse

	if err := c.request(ctx, "getTokenAccountBalance", []interface{}{base58Addr}, &res); err != nil {
		return nil, err
	}

	return &res.Result.Value, nil
}

//func (c *Client) GetTokenAccountsByOwner(ctx context.Context, mint, programID string) (*rpc.GetTokenAccountsByOwnerResponse)