-- name: GetUserVerificationByUserID :one
SELECT *
FROM user_verifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;
-- name: CreateUserVerification :exec
INSERT INTO user_verifications (user_id, email, verification_code)
VALUES (@user_id, @email, @verification_code) ON CONFLICT (user_id, email) DO
UPDATE
SET verification_code = @verification_code;
-- name: DeleteUserVerificationsByEmail :exec
DELETE FROM user_verifications
WHERE email = @email;
-- name: DeleteUserVerificationsByUserID :exec
DELETE FROM user_verifications
WHERE user_id = @user_id;