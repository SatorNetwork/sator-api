// Code generated by sqlc. DO NOT EDIT.
// source: challenge.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const getChallengeByID = `-- name: GetChallengeByID :one
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at
FROM challenges
WHERE id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetChallengeByID(ctx context.Context, id uuid.UUID) (Challenge, error) {
	row := q.queryRow(ctx, q.getChallengeByIDStmt, getChallengeByID, id)
	var i Challenge
	err := row.Scan(
		&i.ID,
		&i.ShowID,
		&i.Title,
		&i.Description,
		&i.PrizePool,
		&i.PlayersToStart,
		&i.TimePerQuestion,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getChallenges = `-- name: GetChallenges :many
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at
FROM challenges WHERE show_id = $1
ORDER BY updated_at,
    created_at DESC
LIMIT $2 OFFSET $3
`

type GetChallengesParams struct {
	ShowID uuid.UUID `json:"show_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) GetChallenges(ctx context.Context, arg GetChallengesParams) ([]Challenge, error) {
	rows, err := q.query(ctx, q.getChallengesStmt, getChallenges, arg.ShowID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Challenge
	for rows.Next() {
		var i Challenge
		if err := rows.Scan(
			&i.ID,
			&i.ShowID,
			&i.Title,
			&i.Description,
			&i.PrizePool,
			&i.PlayersToStart,
			&i.TimePerQuestion,
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
