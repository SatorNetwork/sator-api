-- name: GetEpisodeRatingByID :one
SELECT AVG(rating)::FLOAT as avg_rating
FROM ratings
WHERE episode_id = $1 group by episode_id;
-- name: RateEpisode :exec
INSERT INTO ratings (
    episode_id,
    user_id,
    rating
) VALUES (
    @episode_id,
    @user_id,
    @rating
    );