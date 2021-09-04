-- name: GetEpisodeRatingByID :one
SELECT 
    AVG(rating)::FLOAT AS avg_rating,
    COUNT(episode_id) AS ratings
FROM ratings
WHERE episode_id = $1 
GROUP BY episode_id;

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