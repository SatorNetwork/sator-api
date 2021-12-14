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
) ON CONFLICT (episode_id, user_id) DO
UPDATE SET
    rating = EXCLUDED.rating;

-- name: ReviewEpisode :exec
INSERT INTO ratings (
    episode_id,
    user_id,
    username,
    rating,
    title,
    review
) VALUES (
    @episode_id,
    @user_id,
    @username,
    @rating,
    @title,
    @review
) ON CONFLICT (episode_id, user_id) DO 
UPDATE SET 
    rating = EXCLUDED.rating, 
    title = EXCLUDED.title, 
    review = EXCLUDED.review,
    username = EXCLUDED.username;

-- name: DidUserRateEpisode :one
SELECT EXISTS(
    SELECT * FROM ratings 
    WHERE user_id = @user_id 
    AND episode_id = @episode_id
);

-- name: DidUserReviewEpisode :one
SELECT EXISTS(
    SELECT * FROM ratings 
    WHERE user_id = @user_id 
    AND episode_id = @episode_id
    AND review IS NOT NULL
);

-- name: ReviewsList :many
SELECT * FROM ratings
WHERE episode_id = $1
AND title IS NOT NULL
AND review IS NOT NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ReviewsListByUserID :many
SELECT * FROM ratings
WHERE user_id = $1
  AND title IS NOT NULL
  AND review IS NOT NULL
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;

-- name: DeleteReview :exec
DELETE FROM ratings
WHERE id = @id;

-- name: ReviewByUserID :one
SELECT * FROM ratings
WHERE user_id = $1
  AND episode_id = $2
  AND title IS NOT NULL
  AND review IS NOT NULL;