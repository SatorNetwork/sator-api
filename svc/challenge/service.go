package challenge

import (
	"context"
	"fmt"
	"log"

	repository2 "github.com/SatorNetwork/sator-api/svc/challenge/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		cr challengesRepository
	}

	// Challenge struct
	// Fields were rearranged to optimize memory usage.
	Challenge struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		PrizePool   string    `json:"prize_pool"`
		Players     int       `json:"players"`
		Winners     int       `json:"winners"`
		Play        string    `json:"play"`
	}

	challengesRepository interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository2.Challenge, error)
		GetChallenges(ctx context.Context, arg repository2.GetChallengesParams) ([]repository2.Challenge, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(cr challengesRepository) *Service {
	if cr == nil {
		log.Fatalln("challenges repository is not set")
	}

	return &Service{cr: cr}
}

// GetChallengeByID ...
func (s *Service) GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	challenge, err := s.cr.GetChallengeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	return challenge, nil
}

// GetChallenges ...
func (s *Service) GetChallenges(ctx context.Context, limit, offset int32, showID uuid.UUID) (interface{}, error) {
	list, err := s.cr.GetChallenges(ctx, repository2.GetChallengesParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}

	return list, nil
}
