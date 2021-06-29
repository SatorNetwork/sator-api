// Code generated by sqlc. DO NOT EDIT.
// source: episodes.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addEpisode = `-- name: AddEpisode :exec
INSERT INTO episodes (
    show_id,
    episode_number,
    cover,
    title,
    description,
    release_date
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
           $6
       )
`

type AddEpisodeParams struct {
	ShowID        uuid.UUID      `json:"show_id"`
	EpisodeNumber int32          `json:"episode_number"`
	Cover         sql.NullString `json:"cover"`
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	ReleaseDate   sql.NullTime   `json:"release_date"`
}

func (q *Queries) AddEpisode(ctx context.Context, arg AddEpisodeParams) error {
	_, err := q.exec(ctx, q.addEpisodeStmt, addEpisode,
		arg.ShowID,
		arg.EpisodeNumber,
		arg.Cover,
		arg.Title,
		arg.Description,
		arg.ReleaseDate,
	)
	return err
}

const deleteEpisodeByID = `-- name: DeleteEpisodeByID :exec
DELETE FROM episodes
WHERE id = $1
`

func (q *Queries) DeleteEpisodeByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteEpisodeByIDStmt, deleteEpisodeByID, id)
	return err
}

const deleteEpisodeByShowID = `-- name: DeleteEpisodeByShowID :exec
DELETE FROM episodes
WHERE show_id = $1
`

func (q *Queries) DeleteEpisodeByShowID(ctx context.Context, showID uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteEpisodeByShowIDStmt, deleteEpisodeByShowID, showID)
	return err
}

const getEpisodeByID = `-- name: GetEpisodeByID :one
SELECT id, show_id, episode_number, cover, title, description, release_date, updated_at, created_at
FROM episodes
WHERE id = $1
`

func (q *Queries) GetEpisodeByID(ctx context.Context, id uuid.UUID) (Episode, error) {
	row := q.queryRow(ctx, q.getEpisodeByIDStmt, getEpisodeByID, id)
	var i Episode
	err := row.Scan(
		&i.ID,
		&i.ShowID,
		&i.EpisodeNumber,
		&i.Cover,
		&i.Title,
		&i.Description,
		&i.ReleaseDate,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getEpisodesByShowID = `-- name: GetEpisodesByShowID :many
SELECT id, show_id, episode_number, cover, title, description, release_date, updated_at, created_at
FROM episodes
WHERE show_id = $1
ORDER BY episode_number DESC
    LIMIT $2 OFFSET $3
`

type GetEpisodesByShowIDParams struct {
	ShowID uuid.UUID `json:"show_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) GetEpisodesByShowID(ctx context.Context, arg GetEpisodesByShowIDParams) ([]Episode, error) {
	rows, err := q.query(ctx, q.getEpisodesByShowIDStmt, getEpisodesByShowID, arg.ShowID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Episode
	for rows.Next() {
		var i Episode
		if err := rows.Scan(
			&i.ID,
			&i.ShowID,
			&i.EpisodeNumber,
			&i.Cover,
			&i.Title,
			&i.Description,
			&i.ReleaseDate,
			&i.UpdatedAt,
			&i.CreatedAt,
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
SET show_id = $1,
    episode_number = $2,
    cover = $3,
    title = $4,
    description = $5,
    release_date = $6
WHERE id = $7
`

type UpdateEpisodeParams struct {
	ShowID        uuid.UUID      `json:"show_id"`
	EpisodeNumber int32          `json:"episode_number"`
	Cover         sql.NullString `json:"cover"`
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	ReleaseDate   sql.NullTime   `json:"release_date"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateEpisode(ctx context.Context, arg UpdateEpisodeParams) error {
	_, err := q.exec(ctx, q.updateEpisodeStmt, updateEpisode,
		arg.ShowID,
		arg.EpisodeNumber,
		arg.Cover,
		arg.Title,
		arg.Description,
		arg.ReleaseDate,
		arg.ID,
	)
	return err
}
