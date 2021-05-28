package quiz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		repo           quizRepository
		tokenGenFunc   tokenGenFunc
		tokenParseFunc tokenParseFunc
		tokenTTL       int64
	}

	quizRepository interface{}

	tokenGenFunc   func(data interface{}, ttl int64) (string, error)
	tokenParseFunc func(token string) (interface{}, error)

	TokenPayload struct {
		UserID          string
		Username        string
		ChallengeRoomID string
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo quizRepository, gfn tokenGenFunc, pfn tokenParseFunc) *Service {
	return &Service{repo: repo, tokenGenFunc: gfn}
}

// GetQuizLink returns link with token to connect to quiz
func (s *Service) GetQuizLink(_ context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (interface{}, error) {
	token, err := s.tokenGenFunc(TokenPayload{
		UserID:          uid.String(),
		Username:        username,
		ChallengeRoomID: uuid.New().String(),
	}, s.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("could not generate new token to connect quiz: %w", err)
	}
	return token, nil
}

// ParseQuizToken returns data from quiz connect token
func (s *Service) ParseQuizToken(_ context.Context, token string) (*TokenPayload, error) {
	payload, err := s.tokenParseFunc(token)
	if err != nil {
		return nil, fmt.Errorf("could not parse connection token: %w", err)
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("could not parse connection token: %w", err)
	}
	result := &TokenPayload{}
	if err := json.Unmarshal(b, result); err != nil {
		return nil, fmt.Errorf("could not parse connection token: %w", err)
	}
	return result, nil
}
