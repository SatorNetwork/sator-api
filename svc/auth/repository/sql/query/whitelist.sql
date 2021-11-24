-- name: IsEmailWhitelisted :one
SELECT count(*) > 0 FROM whitelist
WHERE (allowed_type = 'email_domain'
AND @email::text LIKE CONCAT('%', allowed_value))
OR (allowed_type = 'email'
AND allowed_value = @email::text)
LIMIT 1;

-- name: DeleteFromWhitelist :exec
DELETE FROM whitelist
WHERE allowed_type = $1 AND allowed_value = $2;

-- name: GetWhitelist :many
SELECT *
FROM whitelist
ORDER BY allowed_value ASC
    LIMIT $1 OFFSET $2;

-- name: AddToWhitelist :one
INSERT INTO whitelist (
    allowed_type,
    allowed_value
)
VALUES (
           @allowed_type,
           @allowed_value
       ) RETURNING *;