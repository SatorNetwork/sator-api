-- name: GetUserVerificationByUserID :one
SELECT *
FROM user_verifications
WHERE request_type = @request_type
    AND user_id = @user_id
ORDER BY created_at DESC
LIMIT 1;
-- name: GetUserVerificationByEmail :one
SELECT *
FROM user_verifications
WHERE request_type = @request_type
    AND email = @email
ORDER BY created_at DESC
LIMIT 1;
-- name: CreateUserVerification :exec
INSERT INTO user_verifications (request_type, user_id, email, verification_code)
VALUES (
        @request_type,
        @user_id,
        @email,
        @verification_code
    ) ON CONFLICT (request_type, user_id, email) DO
UPDATE
SET verification_code = @verification_code;
-- name: DeleteUserVerificationsByEmail :exec
DELETE FROM user_verifications
WHERE request_type = @request_type
    AND email = @email;
-- name: DeleteUserVerificationsByUserID :exec
DELETE FROM user_verifications
WHERE request_type = @request_type
    AND user_id = @user_id;