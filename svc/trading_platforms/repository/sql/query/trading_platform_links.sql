-- name: CreateTradingPlatformLink :one
INSERT INTO trading_platform_links (
   title,
   link,
   logo
)
VALUES (
   @title,
   @link,
   @logo
) RETURNING *;

-- name: UpdateTradingPlatformLink :exec
UPDATE trading_platform_links
SET
   title = @title,
   link = @link,
   logo = @logo
WHERE id = @id;

-- name: DeleteTradingPlatformLink :exec
DELETE FROM trading_platform_links
WHERE id = @id;

-- name: GetTradingPlatformLinks :many
SELECT *
FROM trading_platform_links
LIMIT $1 OFFSET $2;
