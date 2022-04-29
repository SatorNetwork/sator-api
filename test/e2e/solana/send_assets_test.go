package solana

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"

	lib_coingecko "github.com/SatorNetwork/sator-api/lib/coingecko"
	"github.com/SatorNetwork/sator-api/lib/fee_accumulator"
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/mock"
)

const (
	defaultFee = 1e6
)

func TestSendAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coingeckoMock := lib_coingecko.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.CoingeckoProvider, coingeckoMock)
	solanaPriceInUSD := float64(100)
	satorPriceInUSD := float64(2)
	simplePriceCallback := func(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
		priceMap := map[string]map[string]float32{
			"solana":  {"usd": float32(solanaPriceInUSD)},
			"sator":   {"usd": float32(satorPriceInUSD)},
			"arweave": {"usd": 1},
		}
		return &priceMap, nil
	}
	coingeckoMock.EXPECT().
		SimplePrice([]string{"solana", "sator", "arweave"}, []string{"usd"}).
		DoAndReturn(simplePriceCallback).
		Times(1)

	c := client.NewClient()

	exchangeRatesClient, err := exchange_rates_client.Easy(c.DB.Client())
	require.NoError(t, err)
	_, err = exchangeRatesClient.SyncExchangeRates(context.Background(), &exchange_rates_svc.Empty{})
	require.NoError(t, err)

	solanaClient := solana_client.New(app_config.AppConfigForTests.SolanaApiBaseUrl, solana_client.Config{
		SystemProgram:         app_config.AppConfigForTests.SolanaSystemProgram,
		SysvarRent:            app_config.AppConfigForTests.SolanaSysvarRent,
		SysvarClock:           app_config.AppConfigForTests.SolanaSysvarClock,
		SplToken:              app_config.AppConfigForTests.SolanaSplToken,
		StakeProgramID:        app_config.AppConfigForTests.SolanaStakeProgramID,
		TokenHolderAddr:       app_config.AppConfigForTests.SolanaTokenHolderAddr,
		FeeAccumulatorAddress: app_config.AppConfigForTests.FeeAccumulatorAddress,
	}, exchangeRatesClient)

	solanaFeePayerPrivateKeyBytes, err := base64.StdEncoding.DecodeString(app_config.AppConfigForTests.SolanaFeePayerPrivateKey)
	require.NoError(t, err)
	feePayer, err := types.AccountFromBytes(solanaFeePayerPrivateKeyBytes)
	require.NoError(t, err)
	source := types.NewAccount()
	recipient := types.NewAccount()

	resp, err := solanaClient.PrepareSendAssetsTx(
		context.Background(),
		app_config.AppConfigForTests.SolanaAssetAddr,
		feePayer,
		source,
		recipient.PublicKey.ToBase58(),
		100,
		&lib_solana.SendAssetsConfig{
			PercentToCharge:           5,
			ChargeSolanaFeeFromSender: false,
			AllowFallbackToDefaultFee: true,
			DefaultFee:                defaultFee,
		},
	)
	require.NoError(t, err)
	require.Equal(t, float64(5), resp.FeeInSAO)

	resp, err = solanaClient.PrepareSendAssetsTx(
		context.Background(),
		app_config.AppConfigForTests.SolanaAssetAddr,
		feePayer,
		source,
		recipient.PublicKey.ToBase58(),
		100,
		&lib_solana.SendAssetsConfig{
			PercentToCharge:           5,
			ChargeSolanaFeeFromSender: true,
			AllowFallbackToDefaultFee: true,
			DefaultFee:                defaultFee,
		},
	)
	require.NoError(t, err)
	require.Equal(t, 5+float64(resp.BlockchainFeeInSOLMltpl)/fee_accumulator.SolMltpl*solanaPriceInUSD/satorPriceInUSD, resp.FeeInSAO)
}
