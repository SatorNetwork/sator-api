package challenge

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SatorNetwork/sator-api/svc/challenge/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		cr        challengesRepository
		playUrlFn playURLGenerator
	}

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	// Challenge struct
	// Fields were rearranged to optimize memory usage.
	Challenge struct {
		ID                 uuid.UUID `json:"id"`
		ShowID             uuid.UUID `json:"show_id"`
		Title              string    `json:"title"`
		Description        string    `json:"description"`
		PrizePool          string    `json:"prize_pool"`
		PrizePoolAmount    float64   `json:"-"`
		Players            int       `json:"players"`
		Winners            string    `json:"winners"`
		TimePerQuestion    string    `json:"time_per_question"`
		TimePerQuestionSec int64     `json:"-"`
		Play               string    `json:"play"`
	}

	challengesRepository interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository.Challenge, error)
		GetChallenges(ctx context.Context, arg repository.GetChallengesParams) ([]repository.Challenge, error)
		AddChallenge(ctx context.Context, arg repository.AddChallengeParams) error
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, arg repository.UpdateChallengeParams) error
	}

	playURLGenerator func(challengeID uuid.UUID) string
)

// DefaultPlayURLGenerator ...
func DefaultPlayURLGenerator(baseURL string) playURLGenerator {
	return func(challengeID uuid.UUID) string {
		return fmt.Sprintf("%s/%s/play", baseURL, challengeID.String())
	}
}

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(cr challengesRepository, fn playURLGenerator) *Service {
	if cr == nil {
		log.Fatalln("challenges repository is not set")
	}

	return &Service{
		cr:        cr,
		playUrlFn: fn,
	}
}

// GetByID ...
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (Challenge, error) {
	challenge, err := s.cr.GetChallengeByID(ctx, id)
	if err != nil {
		return Challenge{}, fmt.Errorf("could not get challenge by id: %w", err)
	}

	return castToChallenge(challenge, s.playUrlFn), nil
}

// GetChallengeByID ...
func (s *Service) GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	return s.GetByID(ctx, id)
}

// GetChallengesByShowID ...
func (s *Service) GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	list, err := s.cr.GetChallenges(ctx, repository.GetChallengesParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}

	// Cast repository.Challenge into challenge.Challenge struct
	result := make([]Challenge, 0, len(list))
	for _, v := range list {
		result = append(result, castToChallenge(v, s.playUrlFn))
	}
	return result, nil
}

func castToChallenge(c repository.Challenge, playUrlFn playURLGenerator) Challenge {
	return Challenge{
		ID:                 c.ID,
		Title:              c.Title,
		Description:        c.Description.String,
		PrizePool:          fmt.Sprintf("%.2f SAO", c.PrizePool),
		PrizePoolAmount:    c.PrizePool,
		Players:            int(c.PlayersToStart),
		TimePerQuestion:    fmt.Sprintf("%d sec", c.TimePerQuestion.Int32),
		TimePerQuestionSec: int64(c.TimePerQuestion.Int32),
		Play:               playUrlFn(c.ID),
	}
}

// AddChallenge ..
func (s *Service) AddChallenge(ctx context.Context, ch Challenge) error {
	timePerQuestion, err := strconv.Atoi(ch.TimePerQuestion)
	if err != nil {
		return fmt.Errorf("can not parse TimePerQuestion value from string")
	}

	prizePool, err := strconv.Atoi(ch.PrizePool)
	if err != nil {
		return fmt.Errorf("can not parse PrizePool value from string")
	}

	if err := s.cr.AddChallenge(ctx, repository.AddChallengeParams{
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Description,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      float64(prizePool),
		PlayersToStart: int32(ch.Players),
		TimePerQuestion: sql.NullInt32{
			Int32: int32(timePerQuestion),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not add challenge with title=%s: %w", ch.Title, err)
	}
	return nil
}

// DeleteChallengeByID ...
func (s *Service) DeleteChallengeByID(ctx context.Context, id uuid.UUID) error {
	if err := s.cr.DeleteChallengeByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete challenge with id=%s:%w", id, err)
	}

	return nil
}

// UpdateChallenge ..
func (s *Service) UpdateChallenge(ctx context.Context, ch Challenge) error {
	timePerQuestion, err := strconv.Atoi(ch.TimePerQuestion)
	if err != nil {
		return fmt.Errorf("can not parse TimePerQuestion value from string")
	}

	prizePool, err := strconv.Atoi(ch.PrizePool)
	if err != nil {
		return fmt.Errorf("can not parse PrizePool value from string")
	}

	if err := s.cr.UpdateChallenge(ctx, repository.UpdateChallengeParams{
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Title,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      float64(prizePool),
		PlayersToStart: int32(ch.Players),
		TimePerQuestion: sql.NullInt32{
			Int32: int32(timePerQuestion),
			Valid: false,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: ch.ID,
	}); err != nil {
		return fmt.Errorf("could not update challenge with id=%s:%w", ch.ID, err)
	}
	return nil
}
