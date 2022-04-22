package client

import (
	"context"

	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
)

type (
	Client struct {
		s service
	}

	service interface {
		GetAssetPrice(ctx context.Context, req *exchange_rates_svc.Asset) (*exchange_rates_svc.Price, error)
	}
)

func New(s service) *Client {
	return &Client{s: s}
}

func (c *Client) GetAssetPrice(ctx context.Context, req *exchange_rates_svc.Asset) (*exchange_rates_svc.Price, error) {
	return c.s.GetAssetPrice(ctx, req)
}
