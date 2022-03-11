-- name: LikeDislikeEpisodeReview :exec
INSERT INTO reviews_rating (
    review_id,
    user_id,
    rating_type
) VALUES (
             @review_id,
             @user_id,
             @rating_type
         ) ON CONFLICT (review_id, user_id) DO
UPDATE SET
    rating_type = EXCLUDED.rating_type;

-- name: IsUserRatedReview :one
SELECT count(*) > 0 
FROM reviews_rating
WHERE user_id = @user_id 
AND review_id = @review_id
AND rating_type = @rating_type;

-- name: GetReviewRating :one
SELECT count(*)
FROM reviews_rating
WHERE review_id = @review_id
AND rating_type = @rating_type;

-- name: GetUserEpisodeReview :one
SELECT * FROM reviews_rating
WHERE user_id = @user_id 
AND review_id = @review_id;

-- name: DeleteUserEpisodeReview :exec
DELETE FROM reviews_rating
WHERE user_id = @user_id 
AND review_id = @review_id;
