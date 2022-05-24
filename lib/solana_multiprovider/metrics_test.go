package solana_multiprovider

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	metrics_repository "github.com/SatorNetwork/sator-api/svc/metrics/repository"
	"github.com/SatorNetwork/sator-api/test/framework/client"
)

func TestMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := client.NewClient()
	ctxb := context.Background()
	metricsRepository, err := metrics_repository.Prepare(ctxb, c.DB.Client())
	require.NoError(t, err)

	provider1 := uuid.New().String()
	m := newMetricsRegistrator(metricsRepository)
	{
		m.registerNotAvailableError(ctxb, provider1)
		m.registerOtherError(ctxb, provider1)
		m.registerSuccessCall(ctxb, provider1)

		solanaMetric, err := metricsRepository.GetProviderMetricByName(ctxb, provider1)
		require.NoError(t, err)
		require.Equal(t, provider1, solanaMetric.ProviderName)
		require.Equal(t, int32(1), solanaMetric.NotAvailableErrors)
		require.Equal(t, int32(1), solanaMetric.OtherErrors)
		require.Equal(t, int32(1), solanaMetric.SuccessCalls)

		m.registerNotAvailableError(ctxb, provider1)
		m.registerOtherError(ctxb, provider1)
		m.registerSuccessCall(ctxb, provider1)

		solanaMetric, err = metricsRepository.GetProviderMetricByName(ctxb, provider1)
		require.NoError(t, err)
		require.Equal(t, provider1, solanaMetric.ProviderName)
		require.Equal(t, int32(2), solanaMetric.NotAvailableErrors)
		require.Equal(t, int32(2), solanaMetric.OtherErrors)
		require.Equal(t, int32(2), solanaMetric.SuccessCalls)
	}

	{
		error1 := uuid.New().String()
		error2 := uuid.New().String()
		m.registerError(ctxb, provider1, error1)
		m.registerError(ctxb, provider1, error1)
		m.registerError(ctxb, provider1, error2)

		solanaError, err := m.mr.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
			ProviderName: provider1,
			ErrorMessage: error1,
		})
		require.NoError(t, err)
		require.Equal(t, int32(2), solanaError.Counter)

		solanaError, err = m.mr.GetErrorCounter(ctxb, metrics_repository.GetErrorCounterParams{
			ProviderName: provider1,
			ErrorMessage: error2,
		})
		require.NoError(t, err)
		require.Equal(t, int32(1), solanaError.Counter)
	}
}
