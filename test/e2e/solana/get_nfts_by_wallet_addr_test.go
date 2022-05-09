package solana

import (
	"context"
	"strings"
	"testing"

	"github.com/portto/solana-go-sdk/rpc"
	"github.com/stretchr/testify/require"

	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
)

func TestGetNFTsByWalletAddr(t *testing.T) {
	c := client.NewClient()

	exchangeRatesClient, err := exchange_rates_client.Easy(c.DB.Client())
	require.NoError(t, err)
	_, err = exchangeRatesClient.SyncExchangeRates(context.Background(), &exchange_rates_svc.Empty{})
	require.NoError(t, err)

	solanaClient := solana_client.New(rpc.DevnetRPCEndpoint, solana_client.Config{
		SystemProgram:         app_config.AppConfigForTests.SolanaSystemProgram,
		SysvarRent:            app_config.AppConfigForTests.SolanaSysvarRent,
		SysvarClock:           app_config.AppConfigForTests.SolanaSysvarClock,
		SplToken:              app_config.AppConfigForTests.SolanaSplToken,
		StakeProgramID:        app_config.AppConfigForTests.SolanaStakeProgramID,
		TokenHolderAddr:       app_config.AppConfigForTests.SolanaTokenHolderAddr,
		FeeAccumulatorAddress: app_config.AppConfigForTests.FeeAccumulatorAddress,
	}, exchangeRatesClient)

	ctxb := context.Background()
	addr := "9Qkac1Cyd3bZJ3Hby9N2EWw58q9we3DMYmpft6swoxes"

	nfts, err := solanaClient.GetNFTsByWalletAddress(ctxb, addr)
	if err != nil &&
		!strings.Contains(err.Error(), `{"jsonrpc":"2.0","error":{"code":503,"message":"Service unavailable"}`) &&
		!strings.Contains(err.Error(), `dial tcp: lookup www.arweave.net: No address associated with hostname`) {
		t.Fatal(err)
	}
	_ = nfts
}
