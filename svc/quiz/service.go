package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		repo           quizRepository
		tokenGenFunc   tokenGenFunc
		tokenParseFunc tokenParseFunc
		tokenTTL       int64
		baseQuizURL    string

		quizzes   map[string]broadcast.Broadcaster
		startQuiz chan string // receives quiz id to start
		stopQuiz  chan string // receives quiz id to stop
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
func NewService(repo quizRepository, gfn tokenGenFunc, pfn tokenParseFunc, ttl int64, baseQuizURL string) *Service {
	return &Service{
		repo:           repo,
		tokenGenFunc:   gfn,
		tokenParseFunc: pfn,
		tokenTTL:       ttl,
		baseQuizURL:    baseQuizURL,
	}
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
	return fmt.Sprintf("%s/%s/play/%s", s.baseQuizURL, challengeID, token), nil
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

func (s *Service) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("quiz service stopped")
			return nil
		case qid := <-s.startQuiz:
			go func(c context.Context, id string) {
				if err := s.runQuiz(c, id); err != nil {
					log.Printf("quiz with id %s is finished with error: %v", id, err)
				}
			}(ctx, qid)
		case qid := <-s.stopQuiz:
			if b, ok := s.quizzes[qid]; ok {
				b.Close()
			}
		}
	}
}

func (s *Service) runQuiz(ctx context.Context, qid string) error {

	return nil
}
