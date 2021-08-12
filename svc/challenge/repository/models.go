// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AnswerOption struct {
	ID           uuid.UUID    `json:"id"`
	QuestionID   uuid.UUID    `json:"question_id"`
	AnswerOption string       `json:"answer_option"`
	IsCorrect    sql.NullBool `json:"is_correct"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Attempt struct {
	UserID     uuid.UUID    `json:"user_id"`
	EpisodeID  uuid.UUID    `json:"episode_id"`
	QuestionID uuid.UUID    `json:"question_id"`
	AnswerID   uuid.UUID    `json:"answer_id"`
	Valid      sql.NullBool `json:"valid"`
	CreatedAt  sql.NullTime `json:"created_at"`
}

type Challenge struct {
	ID              uuid.UUID      `json:"id"`
	ShowID          uuid.UUID      `json:"show_id"`
	Title           string         `json:"title"`
	Description     sql.NullString `json:"description"`
	PrizePool       float64        `json:"prize_pool"`
	PlayersToStart  int32          `json:"players_to_start"`
	TimePerQuestion sql.NullInt32  `json:"time_per_question"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	EpisodeID       uuid.UUID      `json:"episode_id"`
	Kind            int32          `json:"kind"`
}

type EpisodeAccess struct {
	EpisodeID   uuid.UUID    `json:"episode_id"`
	UserID      uuid.UUID    `json:"user_id"`
	ActivatedAt sql.NullTime `json:"activated_at"`
}

type Question struct {
	ID            uuid.UUID    `json:"id"`
	ChallengeID   uuid.UUID    `json:"challenge_id"`
	Question      string       `json:"question"`
	QuestionOrder int32        `json:"question_order"`
	UpdatedAt     sql.NullTime `json:"updated_at"`
	CreatedAt     time.Time    `json:"created_at"`
}
