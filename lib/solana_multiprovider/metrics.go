package solana_multiprovider

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"

	metrics_repository "github.com/SatorNetwork/sator-api/svc/metrics/repository"
)

type metricsRepository interface {
	GetProviderMetricByName(ctx context.Context, providerName string) (metrics_repository.SolanaMetric, error)
	UpsertProviderMetrics(ctx context.Context, arg metrics_repository.UpsertProviderMetricsParams) error

	GetErrorCounter(ctx context.Context, arg metrics_repository.GetErrorCounterParams) (metrics_repository.SolanaError, error)
	RegisterProviderError(ctx context.Context, arg metrics_repository.RegisterProviderErrorParams) error
}

type metricsRegistrator struct {
	mr metricsRepository
}

func newMetricsRegistrator(mr metricsRepository) *metricsRegistrator {
	return &metricsRegistrator{
		mr: mr,
	}
}

func (m *metricsRegistrator) registerNotAvailableError(ctx context.Context, providerName string) {
	err := m.applyAndRegister(ctx, providerName, func(m *metrics_repository.SolanaMetric) {
		m.NotAvailableErrors++
	})
	if err != nil {
		log.Println(err)
	}
}

func (m *metricsRegistrator) registerOtherError(ctx context.Context, providerName string) {
	err := m.applyAndRegister(ctx, providerName, func(m *metrics_repository.SolanaMetric) {
		m.OtherErrors++
	})
	if err != nil {
		log.Println(err)
	}
}

func (m *metricsRegistrator) registerSuccessCall(ctx context.Context, providerName string) {
	err := m.applyAndRegister(ctx, providerName, func(m *metrics_repository.SolanaMetric) {
		m.SuccessCalls++
	})
	if err != nil {
		log.Println(err)
	}
}

func (m *metricsRegistrator) applyAndRegister(
	ctx context.Context,
	providerName string,
	apply func(*metrics_repository.SolanaMetric),
) error {
	metric, err := m.mr.GetProviderMetricByName(ctx, providerName)
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "can't get provider metric by name")
	}
	if err != nil && err == sql.ErrNoRows {
		metric = metrics_repository.SolanaMetric{
			ProviderName: providerName,
		}
	}
	apply(&metric)

	err = m.mr.UpsertProviderMetrics(ctx, metrics_repository.UpsertProviderMetricsParams{
		ProviderName:       metric.ProviderName,
		NotAvailableErrors: metric.NotAvailableErrors,
		OtherErrors:        metric.OtherErrors,
		SuccessCalls:       metric.SuccessCalls,
	})
	if err != nil {
		return errors.Wrap(err, "can't upsert provider metrics")
	}

	return nil
}

func (m *metricsRegistrator) registerError(ctx context.Context, providerName, errorMessage string) {
	err := m.mr.RegisterProviderError(ctx, metrics_repository.RegisterProviderErrorParams{
		ProviderName: providerName,
		ErrorMessage: errorMessage,
	})
	if err != nil {
		log.Println(err)
	}
}
