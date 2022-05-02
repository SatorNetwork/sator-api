package fee_accumulator

import (
	"context"

	pkg_errors "github.com/pkg/errors"

	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
)

const (
	SolMltpl = 1e9
	SaoMltpl = 1e9
	ArMltpl  = 1e12
)

type feeAccumulator struct {
	balanceInUSD float64

	solPriceInUSD float64
	saoPriceInUSD float64
	arPriceInUSD  float64
}

func New(exchangeRatesClient *exchange_rates_client.Client) (*feeAccumulator, error) {
	ctx := context.Background()
	solPriceInUSD, err := exchangeRatesClient.GetAssetPrice(ctx, &exchange_rates_svc.Asset{
		AssetType: exchange_rates_svc.AssetTypeSOL,
	})
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't get asset price")
	}
	saoPriceInUSD, err := exchangeRatesClient.GetAssetPrice(ctx, &exchange_rates_svc.Asset{
		AssetType: exchange_rates_svc.AssetTypeSAO,
	})
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't get asset price")
	}
	arPriceInUSD, err := exchangeRatesClient.GetAssetPrice(ctx, &exchange_rates_svc.Asset{
		AssetType: exchange_rates_svc.AssetTypeAR,
	})
	if err != nil {
		return nil, pkg_errors.Wrap(err, "can't get asset price")
	}

	return &feeAccumulator{
		balanceInUSD:  0,
		solPriceInUSD: solPriceInUSD.Usd,
		saoPriceInUSD: saoPriceInUSD.Usd,
		arPriceInUSD:  arPriceInUSD.Usd,
	}, nil
}

func (f *feeAccumulator) AddSOL(solAmt float64) {
	f.balanceInUSD += solAmt * f.solPriceInUSD
}

func (f *feeAccumulator) AddSAO(saoAmt float64) {
	f.balanceInUSD += saoAmt * f.saoPriceInUSD
}

func (f *feeAccumulator) AddAR(arAmt float64) {
	f.balanceInUSD += arAmt * f.arPriceInUSD
}

func (f *feeAccumulator) AddWinstons(winstonsAmt uint64) {
	f.AddAR(float64(winstonsAmt) / ArMltpl)
}

func (f *feeAccumulator) GetFeeInSAO() float64 {
	feeInSAO := f.balanceInUSD / f.saoPriceInUSD
	return feeInSAO
}

func (f *feeAccumulator) GetFeeInSAOMltpl() uint64 {
	return uint64(f.GetFeeInSAO() * SaoMltpl)
}
