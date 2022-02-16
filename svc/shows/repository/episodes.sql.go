// Code generated by sqlc. DO NOT EDIT.
// source: episodes.sql

package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
    verification_challenge_id,
    hint_text,
    watch
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
           $9,
           $10,
           $11
) RETURNING id, show_id, season_id, episode_number, cover, title, description, release_date, updated_at, created_at, challenge_id, verification_challenge_id, hint_text, watch, archived
`

type AddEpisodeParams struct {
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.NullUUID  `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
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
		arg.HintText,
		arg.Watch,
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
		&i.HintText,
		&i.Watch,
		&i.Archived,
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
    episodes.id, episodes.show_id, episodes.season_id, episodes.episode_number, episodes.cover, episodes.title, episodes.description, episodes.release_date, episodes.updated_at, episodes.created_at, episodes.challenge_id, episodes.verification_challenge_id, episodes.hint_text, episodes.watch, episodes.archived, 
    seasons.season_number as season_number
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
WHERE episodes.id = $1 AND episodes.archived = FALSE
`

type GetEpisodeByIDRow struct {
	ID                      uuid.UUID      `json:"id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.NullUUID  `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
	Archived                bool           `json:"archived"`
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
		&i.HintText,
		&i.Watch,
		&i.Archived,
		&i.SeasonNumber,
	)
	return i, err
}

const getEpisodeIDByQuizChallengeID = `-- name: GetEpisodeIDByQuizChallengeID :one
SELECT id
FROM episodes
WHERE challenge_id = $1
`

func (q *Queries) GetEpisodeIDByQuizChallengeID(ctx context.Context, challengeID uuid.NullUUID) (uuid.UUID, error) {
	row := q.queryRow(ctx, q.getEpisodeIDByQuizChallengeIDStmt, getEpisodeIDByQuizChallengeID, challengeID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getEpisodeIDByVerificationChallengeID = `-- name: GetEpisodeIDByVerificationChallengeID :one
SELECT id
FROM episodes
WHERE verification_challenge_id = $1
`

func (q *Queries) GetEpisodeIDByVerificationChallengeID(ctx context.Context, verificationChallengeID uuid.NullUUID) (uuid.UUID, error) {
	row := q.queryRow(ctx, q.getEpisodeIDByVerificationChallengeIDStmt, getEpisodeIDByVerificationChallengeID, verificationChallengeID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
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
    episodes.id, episodes.show_id, episodes.season_id, episodes.episode_number, episodes.cover, episodes.title, episodes.description, episodes.release_date, episodes.updated_at, episodes.created_at, episodes.challenge_id, episodes.verification_challenge_id, episodes.hint_text, episodes.watch, episodes.archived, 
    seasons.season_number as season_number,
    coalesce(avg_ratings.avg_rating, 0) as avg_rating,
    coalesce(avg_ratings.ratings, 0) as ratings
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
LEFT JOIN avg_ratings ON episodes.id = avg_ratings.episode_id
WHERE episodes.show_id = $1
AND episodes.archived = FALSE
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
	SeasonID                uuid.NullUUID  `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
	Archived                bool           `json:"archived"`
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
			&i.HintText,
			&i.Watch,
			&i.Archived,
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

const getListEpisodesByIDs = `-- name: GetListEpisodesByIDs :many
SELECT
    episodes.id, episodes.show_id, episodes.season_id, episodes.episode_number, episodes.cover, episodes.title, episodes.description, episodes.release_date, episodes.updated_at, episodes.created_at, episodes.challenge_id, episodes.verification_challenge_id, episodes.hint_text, episodes.watch, episodes.archived,
    seasons.season_number as season_number,
    shows.title as show_title
FROM episodes
JOIN seasons ON seasons.id = episodes.season_id
JOIN shows ON shows.id = episodes.show_id
WHERE episodes.id = ANY($1::uuid[])
AND episodes.archived = FALSE
ORDER BY episodes.episode_number DESC
`

type GetListEpisodesByIDsRow struct {
	ID                      uuid.UUID      `json:"id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.NullUUID  `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
	Archived                bool           `json:"archived"`
	SeasonNumber            int32          `json:"season_number"`
	ShowTitle               string         `json:"show_title"`
}

func (q *Queries) GetListEpisodesByIDs(ctx context.Context, episodeIds []uuid.UUID) ([]GetListEpisodesByIDsRow, error) {
	rows, err := q.query(ctx, q.getListEpisodesByIDsStmt, getListEpisodesByIDs, pq.Array(episodeIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetListEpisodesByIDsRow
	for rows.Next() {
		var i GetListEpisodesByIDsRow
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
			&i.HintText,
			&i.Watch,
			&i.Archived,
			&i.SeasonNumber,
			&i.ShowTitle,
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
    release_date = $9,
    hint_text = $10,
    watch = $11
WHERE id = $12
`

type UpdateEpisodeParams struct {
	EpisodeNumber           int32          `json:"episode_number"`
	SeasonID                uuid.NullUUID  `json:"season_id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
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
		arg.HintText,
		arg.Watch,
		arg.ID,
	)
	return err
}
