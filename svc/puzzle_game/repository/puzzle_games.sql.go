// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: puzzle_games.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createPuzzleGame = `-- name: CreatePuzzleGame :one
INSERT INTO puzzle_games (
    episode_id,
    prize_pool,
    parts_x,
    parts_y
)
VALUES (
    $1,
    $2,
    $3,
    $3
) RETURNING id, episode_id, prize_pool, parts_x, parts_y, updated_at, created_at
`

type CreatePuzzleGameParams struct {
	EpisodeID uuid.UUID `json:"episode_id"`
	PrizePool float64   `json:"prize_pool"`
	PartsX    int32     `json:"parts_x"`
}

func (q *Queries) CreatePuzzleGame(ctx context.Context, arg CreatePuzzleGameParams) (PuzzleGame, error) {
	row := q.queryRow(ctx, q.createPuzzleGameStmt, createPuzzleGame, arg.EpisodeID, arg.PrizePool, arg.PartsX)
	var i PuzzleGame
	err := row.Scan(
		&i.ID,
		&i.EpisodeID,
		&i.PrizePool,
		&i.PartsX,
		&i.PartsY,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const finishPuzzleGame = `-- name: FinishPuzzleGame :exec
UPDATE puzzle_games_attempts
SET status = 2,
    steps_taken = $1,
    rewards_amount = $2,
    bonus_amount = $3
WHERE puzzle_game_id = $4
AND user_id = $5
AND status = 1
AND image IS NOT NULL
RETURNING id, puzzle_game_id, user_id, status, steps, steps_taken, rewards_amount, image, updated_at, created_at, bonus_amount, tiles
`

type FinishPuzzleGameParams struct {
	StepsTaken    int32     `json:"steps_taken"`
	RewardsAmount float64   `json:"rewards_amount"`
	BonusAmount   float64   `json:"bonus_amount"`
	PuzzleGameID  uuid.UUID `json:"puzzle_game_id"`
	UserID        uuid.UUID `json:"user_id"`
}

func (q *Queries) FinishPuzzleGame(ctx context.Context, arg FinishPuzzleGameParams) error {
	_, err := q.exec(ctx, q.finishPuzzleGameStmt, finishPuzzleGame,
		arg.StepsTaken,
		arg.RewardsAmount,
		arg.BonusAmount,
		arg.PuzzleGameID,
		arg.UserID,
	)
	return err
}

const getPuzzleGameByEpisodeID = `-- name: GetPuzzleGameByEpisodeID :one
SELECT id, episode_id, prize_pool, parts_x, parts_y, updated_at, created_at
FROM puzzle_games
WHERE episode_id = $1
`

func (q *Queries) GetPuzzleGameByEpisodeID(ctx context.Context, episodeID uuid.UUID) (PuzzleGame, error) {
	row := q.queryRow(ctx, q.getPuzzleGameByEpisodeIDStmt, getPuzzleGameByEpisodeID, episodeID)
	var i PuzzleGame
	err := row.Scan(
		&i.ID,
		&i.EpisodeID,
		&i.PrizePool,
		&i.PartsX,
		&i.PartsY,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPuzzleGameByID = `-- name: GetPuzzleGameByID :one
SELECT id, episode_id, prize_pool, parts_x, parts_y, updated_at, created_at
FROM puzzle_games
WHERE id = $1
`

func (q *Queries) GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (PuzzleGame, error) {
	row := q.queryRow(ctx, q.getPuzzleGameByIDStmt, getPuzzleGameByID, id)
	var i PuzzleGame
	err := row.Scan(
		&i.ID,
		&i.EpisodeID,
		&i.PrizePool,
		&i.PartsX,
		&i.PartsY,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPuzzleGameCurrentAttempt = `-- name: GetPuzzleGameCurrentAttempt :one
SELECT id, puzzle_game_id, user_id, status, steps, steps_taken, rewards_amount, image, updated_at, created_at, bonus_amount, tiles
FROM puzzle_games_attempts
WHERE user_id = $1 AND puzzle_game_id = $2 AND status = $3
ORDER BY created_at DESC
LIMIT 1
`

type GetPuzzleGameCurrentAttemptParams struct {
	UserID       uuid.UUID `json:"user_id"`
	PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
	Status       int32     `json:"status"`
}

func (q *Queries) GetPuzzleGameCurrentAttempt(ctx context.Context, arg GetPuzzleGameCurrentAttemptParams) (PuzzleGamesAttempt, error) {
	row := q.queryRow(ctx, q.getPuzzleGameCurrentAttemptStmt, getPuzzleGameCurrentAttempt, arg.UserID, arg.PuzzleGameID, arg.Status)
	var i PuzzleGamesAttempt
	err := row.Scan(
		&i.ID,
		&i.PuzzleGameID,
		&i.UserID,
		&i.Status,
		&i.Steps,
		&i.StepsTaken,
		&i.RewardsAmount,
		&i.Image,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.BonusAmount,
		&i.Tiles,
	)
	return i, err
}

const getPuzzleGameImageIDs = `-- name: GetPuzzleGameImageIDs :many
SELECT file_id FROM puzzle_games_to_images
WHERE puzzle_game_id = $1
`

func (q *Queries) GetPuzzleGameImageIDs(ctx context.Context, puzzleGameID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.query(ctx, q.getPuzzleGameImageIDsStmt, getPuzzleGameImageIDs, puzzleGameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var file_id uuid.UUID
		if err := rows.Scan(&file_id); err != nil {
			return nil, err
		}
		items = append(items, file_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPuzzleGameUnlockOption = `-- name: GetPuzzleGameUnlockOption :one
SELECT id, steps, amount, disabled, updated_at, created_at, locked FROM puzzle_game_unlock_options
WHERE id = $1
`

func (q *Queries) GetPuzzleGameUnlockOption(ctx context.Context, id string) (PuzzleGameUnlockOption, error) {
	row := q.queryRow(ctx, q.getPuzzleGameUnlockOptionStmt, getPuzzleGameUnlockOption, id)
	var i PuzzleGameUnlockOption
	err := row.Scan(
		&i.ID,
		&i.Steps,
		&i.Amount,
		&i.Disabled,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Locked,
	)
	return i, err
}

const getPuzzleGameUnlockOptions = `-- name: GetPuzzleGameUnlockOptions :many
SELECT id, steps, amount, disabled, updated_at, created_at, locked FROM puzzle_game_unlock_options
WHERE disabled = FALSE
ORDER BY locked ASC, amount ASC
`

func (q *Queries) GetPuzzleGameUnlockOptions(ctx context.Context) ([]PuzzleGameUnlockOption, error) {
	rows, err := q.query(ctx, q.getPuzzleGameUnlockOptionsStmt, getPuzzleGameUnlockOptions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PuzzleGameUnlockOption
	for rows.Next() {
		var i PuzzleGameUnlockOption
		if err := rows.Scan(
			&i.ID,
			&i.Steps,
			&i.Amount,
			&i.Disabled,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.Locked,
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

const getUserAvailableSteps = `-- name: GetUserAvailableSteps :one
SELECT coalesce(steps, 0) FROM puzzle_games_attempts
WHERE user_id = $1 AND puzzle_game_id = $2 AND status = 0
ORDER BY created_at DESC
LIMIT 1
`

type GetUserAvailableStepsParams struct {
	UserID       uuid.UUID `json:"user_id"`
	PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
}

func (q *Queries) GetUserAvailableSteps(ctx context.Context, arg GetUserAvailableStepsParams) (int32, error) {
	row := q.queryRow(ctx, q.getUserAvailableStepsStmt, getUserAvailableSteps, arg.UserID, arg.PuzzleGameID)
	var steps int32
	err := row.Scan(&steps)
	return steps, err
}

const linkImageToPuzzleGame = `-- name: LinkImageToPuzzleGame :exec
INSERT INTO puzzle_games_to_images (
    file_id,
    puzzle_game_id
)
VALUES (
    $1,
    $2
)
`

type LinkImageToPuzzleGameParams struct {
	FileID       uuid.UUID `json:"file_id"`
	PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
}

func (q *Queries) LinkImageToPuzzleGame(ctx context.Context, arg LinkImageToPuzzleGameParams) error {
	_, err := q.exec(ctx, q.linkImageToPuzzleGameStmt, linkImageToPuzzleGame, arg.FileID, arg.PuzzleGameID)
	return err
}

const startPuzzleGame = `-- name: StartPuzzleGame :one
UPDATE puzzle_games_attempts
SET status = 1, image = $1, tiles = $2
WHERE puzzle_game_id = $3 
AND user_id = $4
AND status = 0
AND image IS NULL
RETURNING id, puzzle_game_id, user_id, status, steps, steps_taken, rewards_amount, image, updated_at, created_at, bonus_amount, tiles
`

type StartPuzzleGameParams struct {
	Image        sql.NullString `json:"image"`
	Tiles        sql.NullString `json:"tiles"`
	PuzzleGameID uuid.UUID      `json:"puzzle_game_id"`
	UserID       uuid.UUID      `json:"user_id"`
}

func (q *Queries) StartPuzzleGame(ctx context.Context, arg StartPuzzleGameParams) (PuzzleGamesAttempt, error) {
	row := q.queryRow(ctx, q.startPuzzleGameStmt, startPuzzleGame,
		arg.Image,
		arg.Tiles,
		arg.PuzzleGameID,
		arg.UserID,
	)
	var i PuzzleGamesAttempt
	err := row.Scan(
		&i.ID,
		&i.PuzzleGameID,
		&i.UserID,
		&i.Status,
		&i.Steps,
		&i.StepsTaken,
		&i.RewardsAmount,
		&i.Image,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.BonusAmount,
		&i.Tiles,
	)
	return i, err
}

const unlinkImageFromPuzzleGame = `-- name: UnlinkImageFromPuzzleGame :exec
DELETE FROM puzzle_games_to_images
WHERE file_id = $1 AND puzzle_game_id = $2
`

type UnlinkImageFromPuzzleGameParams struct {
	FileID       uuid.UUID `json:"file_id"`
	PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
}

func (q *Queries) UnlinkImageFromPuzzleGame(ctx context.Context, arg UnlinkImageFromPuzzleGameParams) error {
	_, err := q.exec(ctx, q.unlinkImageFromPuzzleGameStmt, unlinkImageFromPuzzleGame, arg.FileID, arg.PuzzleGameID)
	return err
}

const unlockPuzzleGame = `-- name: UnlockPuzzleGame :one
INSERT INTO puzzle_games_attempts (
    puzzle_game_id,
    user_id,
    steps
) 
VALUES (
    $1,
    $2,
    $3
) RETURNING id, puzzle_game_id, user_id, status, steps, steps_taken, rewards_amount, image, updated_at, created_at, bonus_amount, tiles
`

type UnlockPuzzleGameParams struct {
	PuzzleGameID uuid.UUID `json:"puzzle_game_id"`
	UserID       uuid.UUID `json:"user_id"`
	Steps        int32     `json:"steps"`
}

func (q *Queries) UnlockPuzzleGame(ctx context.Context, arg UnlockPuzzleGameParams) (PuzzleGamesAttempt, error) {
	row := q.queryRow(ctx, q.unlockPuzzleGameStmt, unlockPuzzleGame, arg.PuzzleGameID, arg.UserID, arg.Steps)
	var i PuzzleGamesAttempt
	err := row.Scan(
		&i.ID,
		&i.PuzzleGameID,
		&i.UserID,
		&i.Status,
		&i.Steps,
		&i.StepsTaken,
		&i.RewardsAmount,
		&i.Image,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.BonusAmount,
		&i.Tiles,
	)
	return i, err
}

const updatePuzzleGame = `-- name: UpdatePuzzleGame :one
UPDATE puzzle_games
SET
    prize_pool = $1,
    parts_x = $2,
    parts_y = $2
WHERE id = $3
   RETURNING id, episode_id, prize_pool, parts_x, parts_y, updated_at, created_at
`

type UpdatePuzzleGameParams struct {
	PrizePool float64   `json:"prize_pool"`
	PartsX    int32     `json:"parts_x"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdatePuzzleGame(ctx context.Context, arg UpdatePuzzleGameParams) (PuzzleGame, error) {
	row := q.queryRow(ctx, q.updatePuzzleGameStmt, updatePuzzleGame, arg.PrizePool, arg.PartsX, arg.ID)
	var i PuzzleGame
	err := row.Scan(
		&i.ID,
		&i.EpisodeID,
		&i.PrizePool,
		&i.PartsX,
		&i.PartsY,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updatePuzzleGameAttempt = `-- name: UpdatePuzzleGameAttempt :one
UPDATE puzzle_games_attempts
SET
    status = $1,
    steps = $2,
    steps_taken = $3,
    tiles = $4
WHERE id = $5
    RETURNING id, puzzle_game_id, user_id, status, steps, steps_taken, rewards_amount, image, updated_at, created_at, bonus_amount, tiles
`

type UpdatePuzzleGameAttemptParams struct {
	Status     int32          `json:"status"`
	Steps      int32          `json:"steps"`
	StepsTaken int32          `json:"steps_taken"`
	Tiles      sql.NullString `json:"tiles"`
	ID         uuid.UUID      `json:"id"`
}

func (q *Queries) UpdatePuzzleGameAttempt(ctx context.Context, arg UpdatePuzzleGameAttemptParams) (PuzzleGamesAttempt, error) {
	row := q.queryRow(ctx, q.updatePuzzleGameAttemptStmt, updatePuzzleGameAttempt,
		arg.Status,
		arg.Steps,
		arg.StepsTaken,
		arg.Tiles,
		arg.ID,
	)
	var i PuzzleGamesAttempt
	err := row.Scan(
		&i.ID,
		&i.PuzzleGameID,
		&i.UserID,
		&i.Status,
		&i.Steps,
		&i.StepsTaken,
		&i.RewardsAmount,
		&i.Image,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.BonusAmount,
		&i.Tiles,
	)
	return i, err
}
