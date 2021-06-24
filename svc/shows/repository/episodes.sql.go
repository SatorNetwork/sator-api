// Code generated by sqlc. DO NOT EDIT.
// source: episodes.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

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

const getEpisodes = `-- name: GetEpisodes :many
SELECT id, show_id, episode_number, cover, title, description, release_date, updated_at, created_at
FROM episodes
ORDER BY episode_number DESC
    LIMIT $1 OFFSET $2
`

type GetEpisodesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetEpisodes(ctx context.Context, arg GetEpisodesParams) ([]Episode, error) {
	rows, err := q.query(ctx, q.getEpisodesStmt, getEpisodes, arg.Limit, arg.Offset)
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
