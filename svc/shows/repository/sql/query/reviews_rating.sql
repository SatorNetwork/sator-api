-- name: LikeDislikeEpisodeReview :exec
INSERT INTO reviews_rating (
    review_id,
    user_id,
    like_dislike
) VALUES (
             @review_id,
             @user_id,
             @like_dislike
         ) ON CONFLICT (review_id, user_id) DO
UPDATE SET
    like_dislike = EXCLUDED.like_dislike;