-- name: GetEpisodesByShowID :many
SELECT *
FROM episodes
WHERE show_id = $1
ORDER BY episode_number DESC
    LIMIT $2 OFFSET $3;
-- name: GetEpisodeByID :one
SELECT *
FROM episodes
WHERE id = $1;
-- name: AddEpisode :exec
INSERT INTO episodes (
    show_id,
    episode_number,
    cover,
    title,
    description,
    release_date
)
VALUES (
           @show_id,
           @episode_number,
           @cover,
           @title,
           @description,
           @release_date
       );
-- name: UpdateEpisode :exec
UPDATE episodes
SET show_id = @show_id,
    episode_number = @episode_number,
    cover = @cover,
    title = @title,
    description = @description,
    release_date = @release_date
WHERE id = @id;
-- name: DeleteEpisodeByID :exec
DELETE FROM episodes
WHERE id = @id;