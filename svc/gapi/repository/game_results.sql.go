// Code generated by sqlc. DO NOT EDIT.
// source: game_results.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const finishGame = `-- name: FinishGame :exec
UPDATE unity_game_results SET
    blocks_done = $1,
    finished_at = now()
WHERE id = $2
`

type FinishGameParams struct {
	BlocksDone int32     `json:"blocks_done"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) FinishGame(ctx context.Context, arg FinishGameParams) error {
	_, err := q.exec(ctx, q.finishGameStmt, finishGame, arg.BlocksDone, arg.ID)
	return err
}

const getCurrentGame = `-- name: GetCurrentGame :one
SELECT id, user_id, nft_id, complexity, is_training, blocks_done, finished_at, updated_at, created_at FROM unity_game_results 
WHERE user_id = $1 AND finished_at IS NULL 
ORDER BY created_at DESC LIMIT 1
`

func (q *Queries) GetCurrentGame(ctx context.Context, userID uuid.UUID) (UnityGameResult, error) {
	row := q.queryRow(ctx, q.getCurrentGameStmt, getCurrentGame, userID)
	var i UnityGameResult
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.NFTID,
		&i.Complexity,
		&i.IsTraining,
		&i.BlocksDone,
		&i.FinishedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const startGame = `-- name: StartGame :exec
INSERT INTO unity_game_results (
    user_id,
    nft_id,
    complexity,
    is_training
) VALUES ($1, $2, $3, $4)
`

type StartGameParams struct {
	UserID     uuid.UUID `json:"user_id"`
	NFTID      string    `json:"nft_id"`
	Complexity int32     `json:"complexity"`
	IsTraining bool      `json:"is_training"`
}

func (q *Queries) StartGame(ctx context.Context, arg StartGameParams) error {
	_, err := q.exec(ctx, q.startGameStmt, startGame,
		arg.UserID,
		arg.NFTID,
		arg.Complexity,
		arg.IsTraining,
	)
	return err
}
