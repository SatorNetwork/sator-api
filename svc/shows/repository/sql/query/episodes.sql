-- name: GetPublishedEpisodesByShowID :many
WITH avg_ratings AS (
    SELECT 
        episode_id,
        AVG(rating)::FLOAT AS avg_rating,
        COUNT(episode_id) AS ratings
    FROM ratings
    GROUP BY episode_id
)
SELECT 
    episodes.*, 
    seasons.season_number as season_number,
    coalesce(avg_ratings.avg_rating, 0) as avg_rating,
    coalesce(avg_ratings.ratings, 0) as ratings
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
LEFT JOIN avg_ratings ON episodes.id = avg_ratings.episode_id
WHERE episodes.show_id = $1
AND episodes.status = 'published'::episodes_status_type
ORDER BY episodes.episode_number DESC
    LIMIT $2 OFFSET $3;

-- name: GetAllEpisodesByShowID :many
SELECT 
    episodes.*, 
    seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.show_id = $1
ORDER BY episodes.episode_number DESC
    LIMIT $2 OFFSET $3;

-- name: GetPublishedEpisodeByID :one
SELECT 
    episodes.*, 
    seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.id = $1 AND episodes.status = 'published'::episodes_status_type;

-- name: GetEpisodeByID :one
SELECT 
    episodes.*, 
    seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.id = $1;

-- name: GetPublishedRawEpisodeByID :one
SELECT * FROM episodes
WHERE episodes.id = $1 AND episodes.status = 'published'::episodes_status_type;

-- name: GetRawEpisodeByID :one
SELECT * FROM episodes
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
    challenge_id,
    verification_challenge_id,
    hint_text,
    watch,
    status
)
VALUES (
    @show_id,
    @season_id,
    @episode_number,
    @cover,
    @title,
    @description,
    @release_date,
    @challenge_id,
    @verification_challenge_id,
    @hint_text,
    @watch,
    @status::episodes_status_type
) RETURNING *;

-- name: UpdateEpisode :exec
UPDATE episodes
SET episode_number = @episode_number,
    season_id = @season_id,
    show_id = @show_id,
    challenge_id = @challenge_id,
    verification_challenge_id = @verification_challenge_id,
    cover = @cover,
    title = @title,
    description = @description,
    release_date = @release_date,
    hint_text = @hint_text,
    watch = @watch,
    status = @status::episodes_status_type
WHERE id = @id;

-- name: LinkEpisodeChallenges :exec
UPDATE episodes
SET challenge_id = @challenge_id,
    verification_challenge_id = @verification_challenge_id
WHERE id = @id;

-- name: DeleteEpisodeByID :exec
UPDATE episodes
SET status = 'archived'::episodes_status_type
WHERE id = @id;

-- name: GetEpisodeIDByVerificationChallengeID :one
SELECT id
FROM episodes
WHERE verification_challenge_id = $1;

-- name: GetEpisodeIDByQuizChallengeID :one
SELECT id
FROM episodes
WHERE challenge_id = $1;

-- name: GetPublishedListEpisodesByIDs :many
SELECT
    episodes.*,
    seasons.season_number as season_number,
    shows.title as show_title
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
JOIN shows ON shows.id = episodes.show_id
WHERE episodes.id = ANY(@episode_ids::uuid[])
AND episodes.status = 'published'::episodes_status_type
ORDER BY episodes.episode_number DESC;

-- name: GetEpisodesByStatus :many
SELECT *
FROM episodes
WHERE status = @status::episodes_status_type
LIMIT @limit_val OFFSET @offset_val;