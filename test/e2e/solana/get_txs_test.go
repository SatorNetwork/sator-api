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

func TestGetTxs(t *testing.T) {
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
	assetAddr := "FBDfbe7CFXHHNzDpNBYf4Evcg5GKrThYNjk4wP2xwjwA"
	addr := "14RZsAAjDVC5hsVKJ7fgsRG4AfPtPTZqgzR6gu1o5G1T"

	txs, err := solanaClient.GetTransactionsWithAutoDerive(ctxb, assetAddr, addr)
	if err != nil && !strings.Contains(err.Error(), `{"jsonrpc":"2.0","error":{"code":503,"message":"Service unavailable"}`) {
		t.Fatal(err)
	}
	_ = txs
}
