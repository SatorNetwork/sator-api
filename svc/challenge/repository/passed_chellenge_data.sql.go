// Code generated by sqlc. DO NOT EDIT.
// source: passed_chellenge_data.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addChallengeAttempt = `-- name: AddChallengeAttempt :one
INSERT INTO passed_challenges_data (user_id, challenge_id)
VALUES ($1, $2) RETURNING user_id, challenge_id, reward_amount
`

type AddChallengeAttemptParams struct {
	UserID      uuid.UUID `json:"user_id"`
	ChallengeID uuid.UUID `json:"challenge_id"`
}

func (q *Queries) AddChallengeAttempt(ctx context.Context, arg AddChallengeAttemptParams) (PassedChallengesDatum, error) {
	row := q.queryRow(ctx, q.addChallengeAttemptStmt, addChallengeAttempt, arg.UserID, arg.ChallengeID)
	var i PassedChallengesDatum
	err := row.Scan(&i.UserID, &i.ChallengeID, &i.RewardAmount)
	return i, err
}

const countPassedChallengeAttempts = `-- name: CountPassedChallengeAttempts :one
SELECT COUNT (*)
FROM passed_challenges_data
WHERE user_id = $1 AND challenge_id = $2
`

type CountPassedChallengeAttemptsParams struct {
	UserID      uuid.UUID `json:"user_id"`
	ChallengeID uuid.UUID `json:"challenge_id"`
}

func (q *Queries) CountPassedChallengeAttempts(ctx context.Context, arg CountPassedChallengeAttemptsParams) (int64, error) {
	row := q.queryRow(ctx, q.countPassedChallengeAttemptsStmt, countPassedChallengeAttempts, arg.UserID, arg.ChallengeID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getChallengeReceivedRewardAmount = `-- name: GetChallengeReceivedRewardAmount :one
SELECT reward_amount
FROM passed_challenges_data
WHERE user_id = $1 AND challenge_id = $2
`

type GetChallengeReceivedRewardAmountParams struct {
	UserID      uuid.UUID `json:"user_id"`
	ChallengeID uuid.UUID `json:"challenge_id"`
}

func (q *Queries) GetChallengeReceivedRewardAmount(ctx context.Context, arg GetChallengeReceivedRewardAmountParams) (float64, error) {
	row := q.queryRow(ctx, q.getChallengeReceivedRewardAmountStmt, getChallengeReceivedRewardAmount, arg.UserID, arg.ChallengeID)
	var reward_amount float64
	err := row.Scan(&reward_amount)
	return reward_amount, err
}

const storeChallengeReceivedRewardAmount = `-- name: StoreChallengeReceivedRewardAmount :exec
UPDATE passed_challenges_data
SET reward_amount = $1
WHERE user_id = $2 AND challenge_id = $3
`

type StoreChallengeReceivedRewardAmountParams struct {
	RewardAmount float64   `json:"reward_amount"`
	UserID       uuid.UUID `json:"user_id"`
	ChallengeID  uuid.UUID `json:"challenge_id"`
}

func (q *Queries) StoreChallengeReceivedRewardAmount(ctx context.Context, arg StoreChallengeReceivedRewardAmountParams) error {
	_, err := q.exec(ctx, q.storeChallengeReceivedRewardAmountStmt, storeChallengeReceivedRewardAmount, arg.RewardAmount, arg.UserID, arg.ChallengeID)
	return err
}