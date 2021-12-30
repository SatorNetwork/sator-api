// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Episode struct {
	ID                      uuid.UUID      `json:"id"`
	ShowID                  uuid.UUID      `json:"show_id"`
	SeasonID                uuid.NullUUID  `json:"season_id"`
	EpisodeNumber           int32          `json:"episode_number"`
	Cover                   sql.NullString `json:"cover"`
	Title                   string         `json:"title"`
	Description             sql.NullString `json:"description"`
	ReleaseDate             sql.NullTime   `json:"release_date"`
	UpdatedAt               sql.NullTime   `json:"updated_at"`
	CreatedAt               time.Time      `json:"created_at"`
	ChallengeID             uuid.NullUUID  `json:"challenge_id"`
	VerificationChallengeID uuid.NullUUID  `json:"verification_challenge_id"`
	HintText                sql.NullString `json:"hint_text"`
	Watch                   sql.NullString `json:"watch"`
	Archived                bool           `json:"archived"`
}

type Rating struct {
	EpisodeID uuid.UUID      `json:"episode_id"`
	UserID    uuid.UUID      `json:"user_id"`
	Rating    int32          `json:"rating"`
	CreatedAt time.Time      `json:"created_at"`
	ID        uuid.UUID      `json:"id"`
	Title     sql.NullString `json:"title"`
	Review    sql.NullString `json:"review"`
	Username  sql.NullString `json:"username"`
}

type ReviewsRating struct {
	ReviewID   uuid.UUID     `json:"review_id"`
	UserID     uuid.UUID     `json:"user_id"`
	RatingType sql.NullInt32 `json:"rating_type"`
	CreatedAt  time.Time     `json:"created_at"`
}

type Season struct {
	ID           uuid.UUID `json:"id"`
	ShowID       uuid.UUID `json:"show_id"`
	SeasonNumber int32     `json:"season_number"`
}

type Show struct {
	ID             uuid.UUID      `json:"id"`
	Title          string         `json:"title"`
	Cover          string         `json:"cover"`
	HasNewEpisode  bool           `json:"has_new_episode"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	Category       sql.NullString `json:"category"`
	Description    sql.NullString `json:"description"`
	RealmsTitle    sql.NullString `json:"realms_title"`
	RealmsSubtitle sql.NullString `json:"realms_subtitle"`
	Watch          sql.NullString `json:"watch"`
	Archived       bool           `json:"archived"`
}

type ShowClap struct {
	ID        uuid.UUID `json:"id"`
	ShowID    uuid.UUID `json:"show_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
