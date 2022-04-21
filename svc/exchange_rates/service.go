package exchange_rates

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"github.com/SatorNetwork/sator-api/lib/coingecko"
	coingecko_client "github.com/SatorNetwork/sator-api/lib/coingecko/client"
	exchange_rates_repository "github.com/SatorNetwork/sator-api/svc/exchange_rates/repository"
)

type AssetType uint8

const (
	AssetTypeSOL AssetType = iota
	AssetTypeSAO
	AssetTypeAR
)

func (a AssetType) CoingeckoAssetID() string {
	switch a {
	case AssetTypeSOL:
		return "solana"
	case AssetTypeSAO:
		return "sator"
	case AssetTypeAR:
		return "arweave"
	default:
		return ""
	}
}

type (
	exchangeRatesServer struct {
		er exchangeRatesRepository

		coingecko coingecko.Interface
	}

	exchangeRatesRepository interface {
		UpsertExchangeRate(ctx context.Context, arg exchange_rates_repository.UpsertExchangeRateParams) error
		GetExchangeRateByAssetType(ctx context.Context, assetType string) (exchange_rates_repository.ExchangeRate, error)
	}

	Asset struct {
		AssetType AssetType `json:"asset_type"`
	}

	Price struct {
		Usd float64 `json:"usd"`
	}

	Empty struct{}
)

func NewExchangeRatesServer(
	er exchangeRatesRepository,
) (*exchangeRatesServer, error) {
	s := &exchangeRatesServer{
		er:        er,
		coingecko: coingecko_client.NewCoingeckoClient(),
	}
	s.start()

	return s, nil
}

func (s *exchangeRatesServer) start() {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		if err := s.syncExchangeRates(); err != nil {
			log.Printf("can't sync exchange rates: %v", err)
		}
	})
	if err != nil {
		log.Printf("can't register sync-exchange-rates callback")
	}

	c.Start()
}

func (s *exchangeRatesServer) syncExchangeRates() error {
	solID := AssetTypeSOL.CoingeckoAssetID()
	saoID := AssetTypeSAO.CoingeckoAssetID()
	arID := AssetTypeAR.CoingeckoAssetID()

	ids := []string{solID, saoID, arID}
	priceMap, err := s.coingecko.SimplePrice(ids, []string{"usd"})
	if err != nil {
		return errors.Wrap(err, "can't get price map from coingecko")
	}
	if priceMap == nil {
		return errors.Errorf("price map is nil")
	}
	if _, ok := (*priceMap)[solID]; !ok {
		return errors.Errorf("price map doesn't contain price for SOL")
	}
	if _, ok := (*priceMap)[saoID]; !ok {
		return errors.Errorf("price map doesn't contain price for SAO")
	}
	if _, ok := (*priceMap)[arID]; !ok {
		return errors.Errorf("price map doesn't contain price for AR")
	}

	solPrice := (*priceMap)[solID]["usd"]
	saoPrice := (*priceMap)[saoID]["usd"]
	arPrice := (*priceMap)[arID]["usd"]
	if solPrice != 0 {
		err := s.er.UpsertExchangeRate(context.Background(), exchange_rates_repository.UpsertExchangeRateParams{
			AssetType: AssetTypeSOL.CoingeckoAssetID(),
			UsdPrice:  float64(solPrice),
		})
		if err != nil {
			return errors.Wrapf(err, "can't add exchange rate for %v", AssetTypeSOL.CoingeckoAssetID())
		}
	}
	if saoPrice != 0 {
		err := s.er.UpsertExchangeRate(context.Background(), exchange_rates_repository.UpsertExchangeRateParams{
			AssetType: AssetTypeSAO.CoingeckoAssetID(),
			UsdPrice:  float64(saoPrice),
		})
		if err != nil {
			return errors.Wrapf(err, "can't add exchange rate for %v", AssetTypeSAO.CoingeckoAssetID())
		}
	}
	if arPrice != 0 {
		err := s.er.UpsertExchangeRate(context.Background(), exchange_rates_repository.UpsertExchangeRateParams{
			AssetType: AssetTypeAR.CoingeckoAssetID(),
			UsdPrice:  float64(arPrice),
		})
		if err != nil {
			return errors.Wrapf(err, "can't add exchange rate for %v", AssetTypeAR.CoingeckoAssetID())
		}
	}

	return nil
}

func (s *exchangeRatesServer) SyncExchangeRates(ctx context.Context, req *Empty) (*Empty, error) {
	err := s.syncExchangeRates()
	if err != nil {
		return nil, errors.Wrap(err, "can't sync exchange rates")
	}

	return &Empty{}, nil
}

func (s *exchangeRatesServer) GetAssetPrice(ctx context.Context, req *Asset) (*Price, error) {
	exchangeRate, err := s.er.GetExchangeRateByAssetType(ctx, req.AssetType.CoingeckoAssetID())
	if err != nil {
		return nil, errors.Wrap(err, "can't get exchange rate by asset type")
	}

	return &Price{
		Usd: exchangeRate.UsdPrice,
	}, nil
}
