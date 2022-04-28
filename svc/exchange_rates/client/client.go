package client

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_repository "github.com/SatorNetwork/sator-api/svc/exchange_rates/repository"
)

type (
	Client struct {
		s service
	}

	service interface {
		SyncExchangeRates(ctx context.Context, req *exchange_rates_svc.Empty) (*exchange_rates_svc.Empty, error)
		GetAssetPrice(ctx context.Context, req *exchange_rates_svc.Asset) (*exchange_rates_svc.Price, error)
	}
)

func New(s service) *Client {
	return &Client{s: s}
}

func (c *Client) SyncExchangeRates(ctx context.Context, req *exchange_rates_svc.Empty) (*exchange_rates_svc.Empty, error) {
	return c.s.SyncExchangeRates(ctx, req)
}

func (c *Client) GetAssetPrice(ctx context.Context, req *exchange_rates_svc.Asset) (*exchange_rates_svc.Price, error) {
	return c.s.GetAssetPrice(ctx, req)
}

func Easy(dbClient *sql.DB) (*Client, error) {
	exchangeRatesRepository, err := exchange_rates_repository.Prepare(context.Background(), dbClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare exchange rates repository")
	}

	exchangeRatesServer, err := exchange_rates.NewExchangeRatesServer(
		exchangeRatesRepository,
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't create exchange rates server")
	}
	exchangeRatesClient := New(exchangeRatesServer)

	return exchangeRatesClient, nil
}
