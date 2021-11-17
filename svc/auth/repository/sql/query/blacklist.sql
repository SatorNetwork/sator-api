-- name: CheckEmail :one
SELECT count(*) > 0 FROM blacklist
WHERE (restricted_type = 'email_domain'
AND @email::text LIKE '%'+restricted_value)
OR (restricted_type = 'email'
AND restricted_value = @email::text)
LIMIT 1;