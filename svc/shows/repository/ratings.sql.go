// Code generated by sqlc. DO NOT EDIT.
// source: ratings.sql

package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const deleteReview = `-- name: DeleteReview :exec
DELETE FROM ratings
WHERE id = $1
`

func (q *Queries) DeleteReview(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteReviewStmt, deleteReview, id)
	return err
}

const didUserRateEpisode = `-- name: DidUserRateEpisode :one
SELECT EXISTS(
    SELECT episode_id, user_id, rating, created_at, id, title, review, username FROM ratings 
    WHERE user_id = $1 
    AND episode_id = $2
)
`

type DidUserRateEpisodeParams struct {
	UserID    uuid.UUID `json:"user_id"`
	EpisodeID uuid.UUID `json:"episode_id"`
}

func (q *Queries) DidUserRateEpisode(ctx context.Context, arg DidUserRateEpisodeParams) (bool, error) {
	row := q.queryRow(ctx, q.didUserRateEpisodeStmt, didUserRateEpisode, arg.UserID, arg.EpisodeID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const didUserReviewEpisode = `-- name: DidUserReviewEpisode :one
SELECT EXISTS(
    SELECT episode_id, user_id, rating, created_at, id, title, review, username FROM ratings 
    WHERE user_id = $1 
    AND episode_id = $2
    AND review IS NOT NULL
)
`

type DidUserReviewEpisodeParams struct {
	UserID    uuid.UUID `json:"user_id"`
	EpisodeID uuid.UUID `json:"episode_id"`
}

func (q *Queries) DidUserReviewEpisode(ctx context.Context, arg DidUserReviewEpisodeParams) (bool, error) {
	row := q.queryRow(ctx, q.didUserReviewEpisodeStmt, didUserReviewEpisode, arg.UserID, arg.EpisodeID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getEpisodeRatingByID = `-- name: GetEpisodeRatingByID :one
SELECT 
    AVG(rating)::FLOAT AS avg_rating,
    COUNT(episode_id) AS ratings
FROM ratings
WHERE episode_id = $1 
GROUP BY episode_id
`

type GetEpisodeRatingByIDRow struct {
	AvgRating float64 `json:"avg_rating"`
	Ratings   int64   `json:"ratings"`
}

func (q *Queries) GetEpisodeRatingByID(ctx context.Context, episodeID uuid.UUID) (GetEpisodeRatingByIDRow, error) {
	row := q.queryRow(ctx, q.getEpisodeRatingByIDStmt, getEpisodeRatingByID, episodeID)
	var i GetEpisodeRatingByIDRow
	err := row.Scan(&i.AvgRating, &i.Ratings)
	return i, err
}

const getReviewByID = `-- name: GetReviewByID :one
SELECT episode_id, user_id, rating, created_at, id, title, review, username FROM ratings
WHERE id = $1
`

func (q *Queries) GetReviewByID(ctx context.Context, id uuid.UUID) (Rating, error) {
	row := q.queryRow(ctx, q.getReviewByIDStmt, getReviewByID, id)
	var i Rating
	err := row.Scan(
		&i.EpisodeID,
		&i.UserID,
		&i.Rating,
		&i.CreatedAt,
		&i.ID,
		&i.Title,
		&i.Review,
		&i.Username,
	)
	return i, err
}

const getUsersEpisodeRatingByID = `-- name: GetUsersEpisodeRatingByID :one
SELECT rating FROM ratings
WHERE episode_id = $1
  AND user_id = $2
`

type GetUsersEpisodeRatingByIDParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) GetUsersEpisodeRatingByID(ctx context.Context, arg GetUsersEpisodeRatingByIDParams) (int32, error) {
	row := q.queryRow(ctx, q.getUsersEpisodeRatingByIDStmt, getUsersEpisodeRatingByID, arg.EpisodeID, arg.UserID)
	var rating int32
	err := row.Scan(&rating)
	return rating, err
}

const rateEpisode = `-- name: RateEpisode :exec
INSERT INTO ratings (
    episode_id,
    user_id,
    rating
) VALUES (
    $1,
    $2,
    $3
) ON CONFLICT (episode_id, user_id) DO
UPDATE SET
    rating = EXCLUDED.rating
`

type RateEpisodeParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int32     `json:"rating"`
}

func (q *Queries) RateEpisode(ctx context.Context, arg RateEpisodeParams) error {
	_, err := q.exec(ctx, q.rateEpisodeStmt, rateEpisode, arg.EpisodeID, arg.UserID, arg.Rating)
	return err
}

const reviewEpisode = `-- name: ReviewEpisode :exec
INSERT INTO ratings (
    episode_id,
    user_id,
    username,
    rating,
    title,
    review
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) ON CONFLICT (episode_id, user_id) DO 
UPDATE SET 
    rating = EXCLUDED.rating, 
    title = EXCLUDED.title, 
    review = EXCLUDED.review,
    username = EXCLUDED.username
`

type ReviewEpisodeParams struct {
	EpisodeID uuid.UUID      `json:"episode_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Username  sql.NullString `json:"username"`
	Rating    int32          `json:"rating"`
	Title     sql.NullString `json:"title"`
	Review    sql.NullString `json:"review"`
}

func (q *Queries) ReviewEpisode(ctx context.Context, arg ReviewEpisodeParams) error {
	_, err := q.exec(ctx, q.reviewEpisodeStmt, reviewEpisode,
		arg.EpisodeID,
		arg.UserID,
		arg.Username,
		arg.Rating,
		arg.Title,
		arg.Review,
	)
	return err
}

const reviewsList = `-- name: ReviewsList :many
WITH likes_numbers AS (
    SELECT count(*) AS likes_number
    FROM reviews_rating
    WHERE review_id = ratings.id
      AND rating_type = 1
), dislikes_numbers AS (
    SELECT count(*) AS dislikes_number
    FROM reviews_rating
    WHERE review_id = ratings.id
      AND rating_type = 2
)
SELECT ratings.episode_id, ratings.user_id, ratings.rating, ratings.created_at, ratings.id, ratings.title, ratings.review, ratings.username,
       coalesce(likes_numbers.likes_number, 0) as likes_number,
       coalesce(dislikes_numbers.dislikes_number, 0) as dislikes_number
FROM ratings
LEFT JOIN likes_numbers ON ratings.id = reviews_rating.review_id
LEFT JOIN dislikes_numbers ON ratings.id = reviews_rating.review_id
WHERE episode_id = $1
AND title IS NOT NULL
AND review IS NOT NULL
ORDER BY likes_number DESC
LIMIT $2 OFFSET $3
`

type ReviewsListParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type ReviewsListRow struct {
	EpisodeID      uuid.UUID      `json:"episode_id"`
	UserID         uuid.UUID      `json:"user_id"`
	Rating         int32          `json:"rating"`
	CreatedAt      time.Time      `json:"created_at"`
	ID             uuid.UUID      `json:"id"`
	Title          sql.NullString `json:"title"`
	Review         sql.NullString `json:"review"`
	Username       sql.NullString `json:"username"`
	LikesNumber    int64          `json:"likes_number"`
	DislikesNumber int64          `json:"dislikes_number"`
}

func (q *Queries) ReviewsList(ctx context.Context, arg ReviewsListParams) ([]ReviewsListRow, error) {
	rows, err := q.query(ctx, q.reviewsListStmt, reviewsList, arg.EpisodeID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReviewsListRow
	for rows.Next() {
		var i ReviewsListRow
		if err := rows.Scan(
			&i.EpisodeID,
			&i.UserID,
			&i.Rating,
			&i.CreatedAt,
			&i.ID,
			&i.Title,
			&i.Review,
			&i.Username,
			&i.LikesNumber,
			&i.DislikesNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reviewsListByUserID = `-- name: ReviewsListByUserID :many
WITH likes_numbers AS (
    SELECT count(*) AS likes_number
    FROM reviews_rating
    WHERE review_id = ratings.id
      AND rating_type = 1
), dislikes_numbers AS (
    SELECT count(*) AS dislikes_number
    FROM reviews_rating
    WHERE review_id = ratings.id
      AND rating_type = 2
)
SELECT ratings.episode_id, ratings.user_id, ratings.rating, ratings.created_at, ratings.id, ratings.title, ratings.review, ratings.username,
       coalesce(likes_numbers.likes_number, 0) as likes_number,
       coalesce(dislikes_numbers.dislikes_number, 0) as dislikes_number
FROM ratings
         LEFT JOIN likes_numbers ON ratings.id = reviews_rating.review_id
         LEFT JOIN dislikes_numbers ON ratings.id = reviews_rating.review_id
WHERE ratings.user_id = $1
  AND title IS NOT NULL
  AND review IS NOT NULL
ORDER BY likes_number DESC
    LIMIT $2 OFFSET $3
`

type ReviewsListByUserIDParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

type ReviewsListByUserIDRow struct {
	EpisodeID      uuid.UUID      `json:"episode_id"`
	UserID         uuid.UUID      `json:"user_id"`
	Rating         int32          `json:"rating"`
	CreatedAt      time.Time      `json:"created_at"`
	ID             uuid.UUID      `json:"id"`
	Title          sql.NullString `json:"title"`
	Review         sql.NullString `json:"review"`
	Username       sql.NullString `json:"username"`
	LikesNumber    int64          `json:"likes_number"`
	DislikesNumber int64          `json:"dislikes_number"`
}

func (q *Queries) ReviewsListByUserID(ctx context.Context, arg ReviewsListByUserIDParams) ([]ReviewsListByUserIDRow, error) {
	rows, err := q.query(ctx, q.reviewsListByUserIDStmt, reviewsListByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReviewsListByUserIDRow
	for rows.Next() {
		var i ReviewsListByUserIDRow
		if err := rows.Scan(
			&i.EpisodeID,
			&i.UserID,
			&i.Rating,
			&i.CreatedAt,
			&i.ID,
			&i.Title,
			&i.Review,
			&i.Username,
			&i.LikesNumber,
			&i.DislikesNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
