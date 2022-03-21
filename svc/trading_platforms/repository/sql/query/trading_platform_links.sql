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

-- name: UpdateTradingPlatformLink :one
UPDATE trading_platform_links
SET
   title = @title,
   link = @link,
   logo = @logo
WHERE id = @id
RETURNING *;

-- name: DeleteTradingPlatformLink :exec
DELETE FROM trading_platform_links
WHERE id = @id;

-- name: GetTradingPlatformLinks :many
SELECT *
FROM trading_platform_links
ORDER BY title ASC
LIMIT $1 OFFSET $2;
