-- name: GetEpisodeAccessData :one
SELECT *
FROM episode_access
WHERE episode_id = $1 AND user_id = $2 AND activated_at = $3
ORDER BY activated_at DESC
LIMIT 1;
-- name: AddEpisodeAccessData :one
INSERT INTO episode_access (episode_id, user_id, activated_at)
VALUES ($1, $2, $3) RETURNING *;
-- name: DeleteEpisodeAccessData :exec
DELETE FROM episode_access
WHERE episode_id = @episode_id AND user_id = @user_id;
-- name: UpdateEpisodeAccessData :exec
UPDATE episode_access
SET activated_at = @activated_at
WHERE episode_id = @episode_id AND user_id = @user_id;
-- name: DoesUserHaveAccessToEpisode :one
SELECT EXISTS (
    SELECT * 
    FROM episode_access
    WHERE episode_id = @episode_id AND user_id = @user_id AND activated_at > @not_earlier_than
);