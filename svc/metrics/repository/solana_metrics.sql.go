// Code generated by sqlc. DO NOT EDIT.
// source: solana_metrics.sql

package repository

import (
	"context"
)

const getProviderMetricByName = `-- name: GetProviderMetricByName :one
SELECT provider_name, not_available_errors, other_errors, success_calls, updated_at, created_at FROM solana_metrics
WHERE provider_name = $1
`

func (q *Queries) GetProviderMetricByName(ctx context.Context, providerName string) (SolanaMetric, error) {
	row := q.queryRow(ctx, q.getProviderMetricByNameStmt, getProviderMetricByName, providerName)
	var i SolanaMetric
	err := row.Scan(
		&i.ProviderName,
		&i.NotAvailableErrors,
		&i.OtherErrors,
		&i.SuccessCalls,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const upsertProviderMetrics = `-- name: UpsertProviderMetrics :exec
INSERT INTO solana_metrics (
    provider_name,
    not_available_errors,
    other_errors,
    success_calls
)
VALUES (
    $1,
    $2,
    $3,
    $4
) ON CONFLICT (provider_name) DO UPDATE
SET
    not_available_errors = $2,
    other_errors = $3,
    success_calls = $4
`

type UpsertProviderMetricsParams struct {
	ProviderName       string `json:"provider_name"`
	NotAvailableErrors int32  `json:"not_available_errors"`
	OtherErrors        int32  `json:"other_errors"`
	SuccessCalls       int32  `json:"success_calls"`
}

func (q *Queries) UpsertProviderMetrics(ctx context.Context, arg UpsertProviderMetricsParams) error {
	_, err := q.exec(ctx, q.upsertProviderMetricsStmt, upsertProviderMetrics,
		arg.ProviderName,
		arg.NotAvailableErrors,
		arg.OtherErrors,
		arg.SuccessCalls,
	)
	return err
}