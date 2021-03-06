// Code generated by sqlc. DO NOT EDIT.
// source: room_players.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const countPlayersInRoom = `-- name: CountPlayersInRoom :one
SELECT COUNT(user_id) AS players
FROM room_players
WHERE challenge_id = $1
GROUP BY challenge_id
LIMIT 1
`

func (q *Queries) CountPlayersInRoom(ctx context.Context, challengeID uuid.UUID) (int64, error) {
	row := q.queryRow(ctx, q.countPlayersInRoomStmt, countPlayersInRoom, challengeID)
	var players int64
	err := row.Scan(&players)
	return players, err
}

const registerNewPlayer = `-- name: RegisterNewPlayer :exec
INSERT INTO room_players (
    challenge_id,
    user_id
)
VALUES (
    $1,
    $2
) ON CONFLICT (challenge_id, user_id) DO NOTHING
`

type RegisterNewPlayerParams struct {
	ChallengeID uuid.UUID `json:"challenge_id"`
	UserID      uuid.UUID `json:"user_id"`
}

func (q *Queries) RegisterNewPlayer(ctx context.Context, arg RegisterNewPlayerParams) error {
	_, err := q.exec(ctx, q.registerNewPlayerStmt, registerNewPlayer, arg.ChallengeID, arg.UserID)
	return err
}

const unregisterPlayer = `-- name: UnregisterPlayer :exec
DELETE FROM room_players
WHERE challenge_id = $1
  AND user_id = $2
`

type UnregisterPlayerParams struct {
	ChallengeID uuid.UUID `json:"challenge_id"`
	UserID      uuid.UUID `json:"user_id"`
}

func (q *Queries) UnregisterPlayer(ctx context.Context, arg UnregisterPlayerParams) error {
	_, err := q.exec(ctx, q.unregisterPlayerStmt, unregisterPlayer, arg.ChallengeID, arg.UserID)
	return err
}
