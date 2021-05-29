package challenge

import (
	"context"
	"fmt"
	"log"

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
		ID              uuid.UUID `json:"id"`
		Title           string    `json:"title"`
		Description     string    `json:"description"`
		PrizePool       string    `json:"prize_pool"`
		Players         int       `json:"players"`
		Winners         string    `json:"winners"`
		TimePerQuestion string    `json:"time_per_question"`
		Play            string    `json:"play"`
	}

	challengesRepository interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository.Challenge, error)
		GetChallenges(ctx context.Context, arg repository.GetChallengesParams) ([]repository.Challenge, error)
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

// GetChallengeByID ...
func (s *Service) GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	challenge, err := s.cr.GetChallengeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	return castToChallenge(challenge, s.playUrlFn), nil
}

// GetChallengesByShowID ...
// TODO: needs refactoring!
func (s *Service) GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	list, err := s.cr.GetChallenges(ctx, repository.GetChallengesParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}

	// Cast repository.Callange into challenge.Challenge struct
	result := make([]Challenge, 0, len(list))
	for _, v := range list {
		result = append(result, castToChallenge(v, s.playUrlFn))
	}
	return result, nil
}

func castToChallenge(c repository.Challenge, playUrlFn playURLGenerator) Challenge {
	return Challenge{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description.String,
		PrizePool:   fmt.Sprintf("%.2f SAO", c.PrizePool),
		Players:     int(c.PlayersToStart),
		Play:        playUrlFn(c.ID),
	}
}
