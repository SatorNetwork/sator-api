-- name: GetPasswordResetByEmail :one
SELECT *
FROM password_resets
WHERE email = $1
ORDER BY created_at DESC
LIMIT 1;
-- name: CreatePasswordReset :exec
INSERT INTO password_resets (user_id, email, token)
VALUES (@user_id, @email, @token) ON CONFLICT (user_id, email) DO
UPDATE
SET token = @token;
-- name: DeletePasswordResetsByEmail :exec
DELETE FROM password_resets
WHERE email = @email;
-- name: DeletePasswordResetsByUserID :exec
DELETE FROM password_resets
WHERE user_id = @user_id;