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

-- name: GetUsersEpisodeRatingByID :one
SELECT rating FROM ratings
WHERE episode_id = $1
  AND user_id = $2;

-- name: ReviewEpisode :one
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
    username = EXCLUDED.username
RETURNING *;

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
WITH likes_numbers AS (
    SELECT count(*) AS likes_number, review_id
    FROM reviews_rating
    WHERE rating_type = 1
    GROUP BY review_id
), dislikes_numbers AS (
    SELECT count(*) AS dislikes_number, review_id
    FROM reviews_rating
    WHERE rating_type = 2
    GROUP BY review_id
)
SELECT
    ratings.*,
    coalesce(likes_numbers.likes_number, 0) as likes_number,
    coalesce(dislikes_numbers.dislikes_number, 0) as dislikes_number
FROM ratings
    LEFT JOIN likes_numbers ON ratings.id = likes_numbers.review_id
    LEFT JOIN dislikes_numbers ON ratings.id = dislikes_numbers.review_id
WHERE episode_id = $1
AND title IS NOT NULL
AND review IS NOT NULL
ORDER BY likes_number DESC
LIMIT $2 OFFSET $3;

-- name: AllReviewsList :many
SELECT * FROM ratings
LIMIT $1 OFFSET $2;

-- name: GetReviewByID :one
SELECT * FROM ratings
WHERE id = $1;

-- name: ReviewsListByUserID :many
WITH likes_numbers AS (
    SELECT count(*) AS likes_number, review_id
    FROM reviews_rating
    WHERE rating_type = 1
    GROUP BY review_id
), dislikes_numbers AS (
    SELECT count(*) AS dislikes_number, review_id
    FROM reviews_rating
    WHERE rating_type = 2
    GROUP BY review_id
)
SELECT ratings.*,
    coalesce(likes_numbers.likes_number, 0) as likes_number,
    coalesce(dislikes_numbers.dislikes_number, 0) as dislikes_number
FROM ratings
    LEFT JOIN likes_numbers ON ratings.id = likes_numbers.review_id
    LEFT JOIN dislikes_numbers ON ratings.id = dislikes_numbers.review_id
WHERE ratings.user_id = $1
    AND title IS NOT NULL
    AND review IS NOT NULL
ORDER BY likes_number DESC
    LIMIT $2 OFFSET $3;

-- name: DeleteReview :exec
DELETE FROM ratings
WHERE id = @id;
