-- name: GetProfileByUserID :one
SELECT *
FROM profiles
WHERE user_id = $1
LIMIT 1;
-- name: CreateProfile :one
INSERT INTO profiles (user_id, first_name, last_name)
VALUES ($1, $2, $3) RETURNING *;
-- name: UpdateProfileByID :exec
UPDATE profiles
SET first_name = @first_name,
    last_name = @last_name
WHERE id = @id;
-- name: UpdateProfileByUserID :exec
UPDATE profiles
SET first_name = @first_name,
    last_name = @last_name
WHERE user_id = @user_id;
-- name: UpdateAvatar :exec
UPDATE profiles
SET avatar = @avatar
WHERE user_id = @user_id;