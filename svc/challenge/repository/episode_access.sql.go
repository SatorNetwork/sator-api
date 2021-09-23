// Code generated by sqlc. DO NOT EDIT.
// source: episode_access.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addEpisodeAccessData = `-- name: AddEpisodeAccessData :one
INSERT INTO episode_access (episode_id, user_id, activated_at, activated_before)
VALUES ($1, $2, $3, $4) RETURNING episode_id, user_id, activated_at, activated_before
`

type AddEpisodeAccessDataParams struct {
	EpisodeID       uuid.UUID    `json:"episode_id"`
	UserID          uuid.UUID    `json:"user_id"`
	ActivatedAt     sql.NullTime `json:"activated_at"`
	ActivatedBefore sql.NullTime `json:"activated_before"`
}

func (q *Queries) AddEpisodeAccessData(ctx context.Context, arg AddEpisodeAccessDataParams) (EpisodeAccess, error) {
	row := q.queryRow(ctx, q.addEpisodeAccessDataStmt, addEpisodeAccessData,
		arg.EpisodeID,
		arg.UserID,
		arg.ActivatedAt,
		arg.ActivatedBefore,
	)
	var i EpisodeAccess
	err := row.Scan(
		&i.EpisodeID,
		&i.UserID,
		&i.ActivatedAt,
		&i.ActivatedBefore,
	)
	return i, err
}

const deleteEpisodeAccessData = `-- name: DeleteEpisodeAccessData :exec
DELETE FROM episode_access
WHERE episode_id = $1 AND user_id = $2
`

type DeleteEpisodeAccessDataParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteEpisodeAccessData(ctx context.Context, arg DeleteEpisodeAccessDataParams) error {
	_, err := q.exec(ctx, q.deleteEpisodeAccessDataStmt, deleteEpisodeAccessData, arg.EpisodeID, arg.UserID)
	return err
}

const doesUserHaveAccessToEpisode = `-- name: DoesUserHaveAccessToEpisode :one
SELECT EXISTS (
    SELECT episode_id, user_id, activated_at, activated_before 
    FROM episode_access
    WHERE episode_id = $1 AND user_id = $2 AND activated_before > NOW()
)
`

type DoesUserHaveAccessToEpisodeParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) DoesUserHaveAccessToEpisode(ctx context.Context, arg DoesUserHaveAccessToEpisodeParams) (bool, error) {
	row := q.queryRow(ctx, q.doesUserHaveAccessToEpisodeStmt, doesUserHaveAccessToEpisode, arg.EpisodeID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getEpisodeAccessData = `-- name: GetEpisodeAccessData :one
SELECT episode_id, user_id, activated_at, activated_before
FROM episode_access
WHERE episode_id = $1 AND user_id = $2
ORDER BY activated_before DESC, activated_at DESC
LIMIT 1
`

type GetEpisodeAccessDataParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) GetEpisodeAccessData(ctx context.Context, arg GetEpisodeAccessDataParams) (EpisodeAccess, error) {
	row := q.queryRow(ctx, q.getEpisodeAccessDataStmt, getEpisodeAccessData, arg.EpisodeID, arg.UserID)
	var i EpisodeAccess
	err := row.Scan(
		&i.EpisodeID,
		&i.UserID,
		&i.ActivatedAt,
		&i.ActivatedBefore,
	)
	return i, err
}

const listIDsAvailableUserEpisodes = `-- name: ListIDsAvailableUserEpisodes :many
SELECT episode_id
FROM episode_access
WHERE user_id = $1 AND activated_before > NOW()
ORDER BY activated_before DESC, activated_at DESC
    LIMIT $2 OFFSET $3
`

type ListIDsAvailableUserEpisodesParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListIDsAvailableUserEpisodes(ctx context.Context, arg ListIDsAvailableUserEpisodesParams) ([]uuid.UUID, error) {
	rows, err := q.query(ctx, q.listIDsAvailableUserEpisodesStmt, listIDsAvailableUserEpisodes, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var episode_id uuid.UUID
		if err := rows.Scan(&episode_id); err != nil {
			return nil, err
		}
		items = append(items, episode_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const numberUsersWhoHaveAccessToEpisode = `-- name: NumberUsersWhoHaveAccessToEpisode :one
SELECT COUNT(
    EXISTS (
               SELECT episode_id, user_id, activated_at, activated_before
               FROM episode_access
               WHERE episode_id = $1 AND activated_before > NOW()
           )
    )::INT
`

func (q *Queries) NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error) {
	row := q.queryRow(ctx, q.numberUsersWhoHaveAccessToEpisodeStmt, numberUsersWhoHaveAccessToEpisode, episodeID)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const updateEpisodeAccessData = `-- name: UpdateEpisodeAccessData :exec
UPDATE episode_access
SET activated_at = $1, activated_before = $2
WHERE episode_id = $3 AND user_id = $4
`

type UpdateEpisodeAccessDataParams struct {
	ActivatedAt     sql.NullTime `json:"activated_at"`
	ActivatedBefore sql.NullTime `json:"activated_before"`
	EpisodeID       uuid.UUID    `json:"episode_id"`
	UserID          uuid.UUID    `json:"user_id"`
}

func (q *Queries) UpdateEpisodeAccessData(ctx context.Context, arg UpdateEpisodeAccessDataParams) error {
	_, err := q.exec(ctx, q.updateEpisodeAccessDataStmt, updateEpisodeAccessData,
		arg.ActivatedAt,
		arg.ActivatedBefore,
		arg.EpisodeID,
		arg.UserID,
	)
	return err
}
