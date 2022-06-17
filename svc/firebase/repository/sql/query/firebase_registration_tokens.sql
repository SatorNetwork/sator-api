-- name: UpsertRegistrationToken :exec
INSERT INTO firebase_registration_tokens (
    device_id,
    user_id,
    registration_token
)
VALUES (
    @device_id,
    @user_id,
    @registration_token
) ON CONFLICT (device_id) DO UPDATE
SET
    user_id = @user_id,
    registration_token = @registration_token;


-- name: GetRegistrationToken :one
SELECT * FROM firebase_registration_tokens
WHERE device_id = @device_id AND user_id = @user_id;

-- name: DeleteRegistrationToken :exec
DELETE FROM firebase_registration_tokens
WHERE device_id = @device_id AND user_id = @user_id;
