// Code generated by sqlc. DO NOT EDIT.
// source: episodes.sql

package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const addEpisode = `-- name: AddEpisode :one
INSERT INTO episodes (
    show_id,
    season_id,
    episode_number,
    cover,
    title,
    description,
    release_date,
    challenge_id,
    verification_challenge_id
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
           $6,
           $7,
           $8,
           $9
) RETURNING id, show_id, season_id, episode_number, cover, title, description, release_date, updated_at, created_at, challenge_id, verification_challenge_id
`

type AddEpisodeParams struct {
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.UUID      `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	ChallengeID             uuid.UUID      `json:"challenge_id"`
	VerificationChallengeID uuid.UUID      `json:"verification_challenge_id"`
}

func (q *Queries) AddEpisode(ctx context.Context, arg AddEpisodeParams) (Episode, error) {
	row := q.queryRow(ctx, q.addEpisodeStmt, addEpisode,
		arg.ShowID,
		arg.SeasonID,
		arg.EpisodeNumber,
		arg.Cover,
		arg.Title,
		arg.Description,
		arg.ReleaseDate,
		arg.ChallengeID,
		arg.VerificationChallengeID,
	)
	var i Episode
	err := row.Scan(
		&i.ID,
		&i.ShowID,
		&i.SeasonID,
		&i.EpisodeNumber,
		&i.Cover,
		&i.Title,
		&i.Description,
		&i.ReleaseDate,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ChallengeID,
		&i.VerificationChallengeID,
	)
	return i, err
}

const deleteEpisodeByID = `-- name: DeleteEpisodeByID :exec
DELETE FROM episodes
WHERE id = $1
`

func (q *Queries) DeleteEpisodeByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteEpisodeByIDStmt, deleteEpisodeByID, id)
	return err
}

const getEpisodeByID = `-- name: GetEpisodeByID :one
SELECT 
    episodes.id, episodes.show_id, episodes.season_id, episodes.episode_number, episodes.cover, episodes.title, episodes.description, episodes.release_date, episodes.updated_at, episodes.created_at, episodes.challenge_id, episodes.verification_challenge_id, 
    seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.id = $1
`

type GetEpisodeByIDRow struct {
	ID                      uuid.UUID      `json:"id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.UUID      `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.UUID      `json:"challenge_id"`
	VerificationChallengeID uuid.UUID      `json:"verification_challenge_id"`
	SeasonNumber            int32          `json:"season_number"`
}

func (q *Queries) GetEpisodeByID(ctx context.Context, id uuid.UUID) (GetEpisodeByIDRow, error) {
	row := q.queryRow(ctx, q.getEpisodeByIDStmt, getEpisodeByID, id)
	var i GetEpisodeByIDRow
	err := row.Scan(
		&i.ID,
		&i.ShowID,
		&i.SeasonID,
		&i.EpisodeNumber,
		&i.Cover,
		&i.Title,
		&i.Description,
		&i.ReleaseDate,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ChallengeID,
		&i.VerificationChallengeID,
		&i.SeasonNumber,
	)
	return i, err
}

const getEpisodesByShowID = `-- name: GetEpisodesByShowID :many
WITH avg_ratings AS (
    SELECT 
        episode_id,
        AVG(rating)::FLOAT AS avg_rating,
        COUNT(episode_id) AS ratings
    FROM ratings
    GROUP BY episode_id
)
SELECT 
    episodes.id, episodes.show_id, episodes.season_id, episodes.episode_number, episodes.cover, episodes.title, episodes.description, episodes.release_date, episodes.updated_at, episodes.created_at, episodes.challenge_id, episodes.verification_challenge_id, 
    seasons.season_number as season_number,
    avg_ratings.avg_rating as avg_rating,
    avg_ratings.ratings as ratings
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
JOIN avg_ratings ON episodes.id = avg_ratings.episode_id
WHERE episodes.show_id = $1
ORDER BY episodes.episode_number DESC
    LIMIT $2 OFFSET $3
`

type GetEpisodesByShowIDParams struct {
	ShowID uuid.UUID `json:"show_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

type GetEpisodesByShowIDRow struct {
	ID                      uuid.UUID      `json:"id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.UUID      `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.UUID      `json:"challenge_id"`
	VerificationChallengeID uuid.UUID      `json:"verification_challenge_id"`
	SeasonNumber            int32          `json:"season_number"`
	AvgRating               float64        `json:"avg_rating"`
	Ratings                 int64          `json:"ratings"`
}

func (q *Queries) GetEpisodesByShowID(ctx context.Context, arg GetEpisodesByShowIDParams) ([]GetEpisodesByShowIDRow, error) {
	rows, err := q.query(ctx, q.getEpisodesByShowIDStmt, getEpisodesByShowID, arg.ShowID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEpisodesByShowIDRow
	for rows.Next() {
		var i GetEpisodesByShowIDRow
		if err := rows.Scan(
			&i.ID,
			&i.ShowID,
			&i.SeasonID,
			&i.EpisodeNumber,
			&i.Cover,
			&i.Title,
			&i.Description,
			&i.ReleaseDate,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.ChallengeID,
			&i.VerificationChallengeID,
			&i.SeasonNumber,
			&i.AvgRating,
			&i.Ratings,
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

const updateEpisode = `-- name: UpdateEpisode :exec
UPDATE episodes
SET episode_number = $1,
    season_id = $2,
    show_id = $3,
    challenge_id = $4,
    verification_challenge_id = $5,
    cover = $6,
    title = $7,
    description = $8,
    release_date = $9
WHERE id = $10
`

type UpdateEpisodeParams struct {
	EpisodeNumber           int32          `json:"episode_number"`
	SeasonID                uuid.UUID      `json:"season_id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	ChallengeID             uuid.UUID      `json:"challenge_id"`
	VerificationChallengeID uuid.UUID      `json:"verification_challenge_id"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	ID                      uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateEpisode(ctx context.Context, arg UpdateEpisodeParams) error {
	_, err := q.exec(ctx, q.updateEpisodeStmt, updateEpisode,
		arg.EpisodeNumber,
		arg.SeasonID,
		arg.ShowID,
		arg.ChallengeID,
		arg.VerificationChallengeID,
		arg.Cover,
		arg.Title,
		arg.Description,
		arg.ReleaseDate,
		arg.ID,
	)
	return err
}
