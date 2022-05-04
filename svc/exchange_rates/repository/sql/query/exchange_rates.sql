-- name: UpsertExchangeRate :exec
INSERT INTO exchange_rates (
    asset_type,
    usd_price
)
VALUES (
   @asset_type,
   @usd_price
) ON CONFLICT (asset_type) DO
UPDATE
SET usd_price = @usd_price;

-- name: GetExchangeRateByAssetType :one
SELECT *
FROM exchange_rates
WHERE asset_type = $1;
