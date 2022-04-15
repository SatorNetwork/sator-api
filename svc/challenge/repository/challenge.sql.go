// Code generated by sqlc. DO NOT EDIT.
// source: challenge.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addChallenge = `-- name: AddChallenge :one
INSERT INTO challenges (
    show_id,
    title,
    description,
    prize_pool,
    players_to_start,
    time_per_question,
    updated_at,
    episode_id,
    kind,
    user_max_attempts,
    max_winners,
    questions_per_game,
    min_correct_answers,
    percent_for_quiz,
    minimum_reward
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
    $11,
    $12,
    $13,
    $14,
    $15
) RETURNING id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at, episode_id, kind, user_max_attempts, max_winners, questions_per_game, min_correct_answers, percent_for_quiz, minimum_reward
`

type AddChallengeParams struct {
	ShowID            uuid.UUID      `json:"show_id"`
	Title             string         `json:"title"`
	Description       sql.NullString `json:"description"`
	PrizePool         float64        `json:"prize_pool"`
	PlayersToStart    int32          `json:"players_to_start"`
	TimePerQuestion   sql.NullInt32  `json:"time_per_question"`
	UpdatedAt         sql.NullTime   `json:"updated_at"`
	EpisodeID         uuid.NullUUID  `json:"episode_id"`
	Kind              int32          `json:"kind"`
	UserMaxAttempts   int32          `json:"user_max_attempts"`
	MaxWinners        sql.NullInt32  `json:"max_winners"`
	QuestionsPerGame  int32          `json:"questions_per_game"`
	MinCorrectAnswers int32          `json:"min_correct_answers"`
	PercentForQuiz    float64        `json:"percent_for_quiz"`
	MinimumReward     float64        `json:"minimum_reward"`
}

func (q *Queries) AddChallenge(ctx context.Context, arg AddChallengeParams) (Challenge, error) {
	row := q.queryRow(ctx, q.addChallengeStmt, addChallenge,
		arg.ShowID,
		arg.Title,
		arg.Description,
		arg.PrizePool,
		arg.PlayersToStart,
		arg.TimePerQuestion,
		arg.UpdatedAt,
		arg.EpisodeID,
		arg.Kind,
		arg.UserMaxAttempts,
		arg.MaxWinners,
		arg.QuestionsPerGame,
		arg.MinCorrectAnswers,
		arg.PercentForQuiz,
		arg.MinimumReward,
	)
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
		&i.EpisodeID,
		&i.Kind,
		&i.UserMaxAttempts,
		&i.MaxWinners,
		&i.QuestionsPerGame,
		&i.MinCorrectAnswers,
		&i.PercentForQuiz,
		&i.MinimumReward,
	)
	return i, err
}

const deleteChallengeByID = `-- name: DeleteChallengeByID :exec
DELETE FROM challenges
WHERE id = $1
`

func (q *Queries) DeleteChallengeByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteChallengeByIDStmt, deleteChallengeByID, id)
	return err
}

const getChallengeByEpisodeID = `-- name: GetChallengeByEpisodeID :one
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at, episode_id, kind, user_max_attempts, max_winners, questions_per_game, min_correct_answers, percent_for_quiz, minimum_reward
FROM challenges
WHERE episode_id = $1
ORDER BY created_at DESC
    LIMIT 1
`

func (q *Queries) GetChallengeByEpisodeID(ctx context.Context, episodeID uuid.NullUUID) (Challenge, error) {
	row := q.queryRow(ctx, q.getChallengeByEpisodeIDStmt, getChallengeByEpisodeID, episodeID)
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
		&i.EpisodeID,
		&i.Kind,
		&i.UserMaxAttempts,
		&i.MaxWinners,
		&i.QuestionsPerGame,
		&i.MinCorrectAnswers,
		&i.PercentForQuiz,
		&i.MinimumReward,
	)
	return i, err
}

const getChallengeByID = `-- name: GetChallengeByID :one
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at, episode_id, kind, user_max_attempts, max_winners, questions_per_game, min_correct_answers, percent_for_quiz, minimum_reward
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
		&i.EpisodeID,
		&i.Kind,
		&i.UserMaxAttempts,
		&i.MaxWinners,
		&i.QuestionsPerGame,
		&i.MinCorrectAnswers,
		&i.PercentForQuiz,
		&i.MinimumReward,
	)
	return i, err
}

const getChallengeByTitle = `-- name: GetChallengeByTitle :one
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at, episode_id, kind, user_max_attempts, max_winners, questions_per_game, min_correct_answers, percent_for_quiz, minimum_reward
FROM challenges
WHERE title = $1
ORDER BY created_at DESC
    LIMIT 1
`

func (q *Queries) GetChallengeByTitle(ctx context.Context, title string) (Challenge, error) {
	row := q.queryRow(ctx, q.getChallengeByTitleStmt, getChallengeByTitle, title)
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
		&i.EpisodeID,
		&i.Kind,
		&i.UserMaxAttempts,
		&i.MaxWinners,
		&i.QuestionsPerGame,
		&i.MinCorrectAnswers,
		&i.PercentForQuiz,
		&i.MinimumReward,
	)
	return i, err
}

const getChallenges = `-- name: GetChallenges :many
SELECT id, show_id, title, description, prize_pool, players_to_start, time_per_question, updated_at, created_at, episode_id, kind, user_max_attempts, max_winners, questions_per_game, min_correct_answers, percent_for_quiz, minimum_reward
FROM challenges
WHERE show_id = $1
ORDER BY updated_at DESC,
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
			&i.EpisodeID,
			&i.Kind,
			&i.UserMaxAttempts,
			&i.MaxWinners,
			&i.QuestionsPerGame,
			&i.MinCorrectAnswers,
			&i.PercentForQuiz,
			&i.MinimumReward,
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

const updateChallenge = `-- name: UpdateChallenge :exec
UPDATE challenges
SET show_id = $1,
    title = $2,
    description = $3,
    prize_pool = $4,
    players_to_start = $5,
    time_per_question = $6,
    updated_at = $7,
    episode_id = $8,
    kind = $9,
    user_max_attempts = $10,
    max_winners = $11,
    questions_per_game = $12,
    min_correct_answers = $13,
    percent_for_quiz = $14,
    minimum_reward = $15
WHERE id = $16
`

type UpdateChallengeParams struct {
	ShowID            uuid.UUID      `json:"show_id"`
	Title             string         `json:"title"`
	Description       sql.NullString `json:"description"`
	PrizePool         float64        `json:"prize_pool"`
	PlayersToStart    int32          `json:"players_to_start"`
	TimePerQuestion   sql.NullInt32  `json:"time_per_question"`
	UpdatedAt         sql.NullTime   `json:"updated_at"`
	EpisodeID         uuid.NullUUID  `json:"episode_id"`
	Kind              int32          `json:"kind"`
	UserMaxAttempts   int32          `json:"user_max_attempts"`
	MaxWinners        sql.NullInt32  `json:"max_winners"`
	QuestionsPerGame  int32          `json:"questions_per_game"`
	MinCorrectAnswers int32          `json:"min_correct_answers"`
	PercentForQuiz    float64        `json:"percent_for_quiz"`
	MinimumReward     float64        `json:"minimum_reward"`
	ID                uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateChallenge(ctx context.Context, arg UpdateChallengeParams) error {
	_, err := q.exec(ctx, q.updateChallengeStmt, updateChallenge,
		arg.ShowID,
		arg.Title,
		arg.Description,
		arg.PrizePool,
		arg.PlayersToStart,
		arg.TimePerQuestion,
		arg.UpdatedAt,
		arg.EpisodeID,
		arg.Kind,
		arg.UserMaxAttempts,
		arg.MaxWinners,
		arg.QuestionsPerGame,
		arg.MinCorrectAnswers,
		arg.PercentForQuiz,
		arg.MinimumReward,
		arg.ID,
	)
	return err
}
