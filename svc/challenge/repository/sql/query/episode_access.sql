-- name: GetEpisodeAccessData :one
SELECT *
FROM episode_access
WHERE episode_id = $1 AND user_id = $2
ORDER BY activated_before DESC, activated_at DESC
LIMIT 1;
-- name: AddEpisodeAccessData :one
INSERT INTO episode_access (episode_id, user_id, activated_at, activated_before)
VALUES ($1, $2, $3, $4) RETURNING *;
-- name: DeleteEpisodeAccessData :exec
DELETE FROM episode_access
WHERE episode_id = @episode_id AND user_id = @user_id;
-- name: UpdateEpisodeAccessData :exec
UPDATE episode_access
SET activated_at = @activated_at, activated_before = @activated_before
WHERE episode_id = @episode_id AND user_id = @user_id;
-- name: DoesUserHaveAccessToEpisode :one
SELECT EXISTS (
    SELECT * 
    FROM episode_access
    WHERE episode_id = @episode_id AND user_id = @user_id AND activated_before > NOW()
);
-- name: NumberUsersWhoHaveAccessToEpisode :one
SELECT COUNT(*)::INT
FROM episode_access
WHERE episode_id = @episode_id AND activated_before > NOW();
-- name: ListIDsAvailableUserEpisodes :many
SELECT episode_id
FROM episode_access
WHERE user_id = $1 AND activated_before > NOW()
ORDER BY activated_before DESC, activated_at DESC
    LIMIT $2 OFFSET $3;
