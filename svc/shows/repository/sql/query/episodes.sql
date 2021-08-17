-- name: GetEpisodesByShowID :many
SELECT episodes.*, seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.show_id = $1
ORDER BY episodes.episode_number DESC
    LIMIT $2 OFFSET $3;
-- name: GetEpisodeByID :one
SELECT episodes.*, seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.id = $1;
-- name: AddEpisode :one
INSERT INTO episodes (
    show_id,
    season_id,
    episode_number,
    cover,
    title,
    description,
    release_date,
    challenge_id
)
VALUES (
           @show_id,
           @season_id,
           @episode_number,
           @cover,
           @title,
           @description,
           @release_date,
           @challenge_id
) RETURNING *;
-- name: UpdateEpisode :exec
UPDATE episodes
SET episode_number = @episode_number,
    season_id = @season_id,
    show_id = @show_id,
    challenge_id = @challenge_id,
    cover = @cover,
    title = @title,
    description = @description,
    release_date = @release_date
WHERE id = @id;
-- name: DeleteEpisodeByID :exec
DELETE FROM episodes
WHERE id = @id;