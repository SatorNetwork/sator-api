// Code generated by sqlc. DO NOT EDIT.
// source: ratings.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

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

const rateEpisode = `-- name: RateEpisode :exec
INSERT INTO ratings (
    episode_id,
    user_id,
    rating
) VALUES (
    $1,
    $2,
    $3
)
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
SELECT episode_id, user_id, rating, created_at, id, title, review, username FROM ratings 
WHERE episode_id = $1
AND title IS NOT NULL
AND review IS NOT NULL
ORDER BY created_at DESC
`

func (q *Queries) ReviewsList(ctx context.Context, episodeID uuid.UUID) ([]Rating, error) {
	rows, err := q.query(ctx, q.reviewsListStmt, reviewsList, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Rating
	for rows.Next() {
		var i Rating
		if err := rows.Scan(
			&i.EpisodeID,
			&i.UserID,
			&i.Rating,
			&i.CreatedAt,
			&i.ID,
			&i.Title,
			&i.Review,
			&i.Username,
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
