package metrics

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	metrics_repository "github.com/SatorNetwork/sator-api/svc/metrics/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestMetricsDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()

	defer app_config.RunAndWait()()

	c := client.NewClient()
	ctxb := context.Background()
	metricsRepository, err := metrics_repository.Prepare(ctxb, c.DB.Client())
	require.NoError(t, err)

	{
		provider1 := uuid.New().String()
		provider2 := uuid.New().String()
		error1 := uuid.New().String()
		error2 := uuid.New().String()
		err = metricsRepository.RegisterProviderError(ctxb, metrics_repository.RegisterProviderErrorParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		solanaError, err := metricsRepository.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		require.Equal(t, int32(1), solanaError.Counter)

		err = metricsRepository.RegisterProviderError(ctxb, metrics_repository.RegisterProviderErrorParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		solanaError, err = metricsRepository.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		require.Equal(t, int32(2), solanaError.Counter)

		err = metricsRepository.RegisterProviderError(ctxb, metrics_repository.RegisterProviderErrorParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		solanaError, err = metricsRepository.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		require.Equal(t, int32(3), solanaError.Counter)

		{
			err = metricsRepository.RegisterProviderError(ctxb, metrics_repository.RegisterProviderErrorParams{
				ProviderName: provider1,
				ErrorMessage: error2,
			})
			require.NoError(t, err)
			solanaError, err := metricsRepository.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
				ProviderName: provider1,
				ErrorMessage: error2,
			})
			require.NoError(t, err)
			require.Equal(t, int32(1), solanaError.Counter)
		}

		{
			err = metricsRepository.RegisterProviderError(ctxb, metrics_repository.RegisterProviderErrorParams{
				ProviderName: provider2,
				ErrorMessage: error1,
			})
			require.NoError(t, err)
			solanaError, err := metricsRepository.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
				ProviderName: provider2,
				ErrorMessage: error1,
			})
			require.NoError(t, err)
			require.Equal(t, int32(1), solanaError.Counter)
		}
	}
}
