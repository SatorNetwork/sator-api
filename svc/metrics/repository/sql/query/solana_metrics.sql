-- name: UpsertProviderMetrics :exec
INSERT INTO solana_metrics (
    provider_name,
    not_available_errors,
    other_errors,
    success_calls
)
VALUES (
    @provider_name,
    @not_available_errors,
    @other_errors,
    @success_calls
) ON CONFLICT (provider_name) DO UPDATE
SET
    not_available_errors = @not_available_errors,
    other_errors = @other_errors,
    success_calls = @success_calls;

-- name: GetProviderMetricByName :one
SELECT * FROM solana_metrics
WHERE provider_name = @provider_name;
