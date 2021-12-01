-- name: IsEmailBlacklisted :one
SELECT count(*) > 0 FROM blacklist
WHERE (restricted_type = 'email_domain'
AND @email::text LIKE CONCAT('%', restricted_value))
OR (restricted_type = 'email'
AND restricted_value = @email::text)
LIMIT 1;

-- name: DeleteFromBlacklist :exec
DELETE FROM blacklist
WHERE restricted_type = $1 AND restricted_value = $2;

-- name: GetBlacklist :many
SELECT *
FROM blacklist
ORDER BY restricted_value ASC
    LIMIT $1 OFFSET $2;

-- name: AddToBlacklist :one
INSERT INTO blacklist (
    restricted_type,
    restricted_value
)
VALUES (
           @restricted_type,
           @restricted_value
       ) RETURNING *;

-- name: GetBlacklistByRestrictedValue :many
SELECT *
FROM blacklist
WHERE restricted_value LIKE CONCAT('%', @query::text, '%')
ORDER BY restricted_value ASC
    LIMIT @limit_val::INT OFFSET @offset_val::INT;