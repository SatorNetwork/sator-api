-- name: IsEmailWhitelisted :one
SELECT count(*) > 0 FROM whitelist
WHERE (allowed_type = 'email_domain'
AND @email::text LIKE CONCAT('%', allowed_value))
OR (allowed_type = 'email'
AND allowed_value = @email::text)
LIMIT 1;