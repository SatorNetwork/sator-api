// Code generated by sqlc. DO NOT EDIT.
// source: questions.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addQuestion = `-- name: AddQuestion :one
INSERT INTO questions (challenge_id, question, question_order)
VALUES ($1, $2, $3) RETURNING id, challenge_id, question, question_order, updated_at, created_at
`

type AddQuestionParams struct {
	ChallengeID   uuid.UUID `json:"challenge_id"`
	Question      string    `json:"question"`
	QuestionOrder int32     `json:"question_order"`
}

func (q *Queries) AddQuestion(ctx context.Context, arg AddQuestionParams) (Question, error) {
	row := q.queryRow(ctx, q.addQuestionStmt, addQuestion, arg.ChallengeID, arg.Question, arg.QuestionOrder)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.ChallengeID,
		&i.Question,
		&i.QuestionOrder,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getQuestionByID = `-- name: GetQuestionByID :one
SELECT id, challenge_id, question, question_order, updated_at, created_at
FROM questions
WHERE id = $1
    LIMIT 1
`

func (q *Queries) GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error) {
	row := q.queryRow(ctx, q.getQuestionByIDStmt, getQuestionByID, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.ChallengeID,
		&i.Question,
		&i.QuestionOrder,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getQuestionsByChallengeID = `-- name: GetQuestionsByChallengeID :many
SELECT id, challenge_id, question, question_order, updated_at, created_at
FROM questions
WHERE challenge_id = $1
ORDER BY quiestion_order ASC
`

func (q *Queries) GetQuestionsByChallengeID(ctx context.Context, challengeID uuid.UUID) ([]Question, error) {
	rows, err := q.query(ctx, q.getQuestionsByChallengeIDStmt, getQuestionsByChallengeID, challengeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.ChallengeID,
			&i.Question,
			&i.QuestionOrder,
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
