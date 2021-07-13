package challenge

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/SatorNetwork/sator-api/svc/questions"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		cr        challengesRepository
		playUrlFn playURLGenerator
		qs        questionsService
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
		Players            int32     `json:"players"`
		Winners            string    `json:"winners"`
		TimePerQuestion    string    `json:"time_per_question"`
		TimePerQuestionSec int64     `json:"-"`
		Play               string    `json:"play"`
		EpisodeID          uuid.UUID `json:"episode_id"`
		Kind               int32     `json:"kind"`
	}

	challengesRepository interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository.Challenge, error)
		GetChallenges(ctx context.Context, arg repository.GetChallengesParams) ([]repository.Challenge, error)

		GetChallengeByEpisodeID(ctx context.Context, episodeID uuid.UUID) (repository.Challenge, error)

		AddChallenge(ctx context.Context, arg repository.AddChallengeParams) (repository.Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, arg repository.UpdateChallengeParams) error
	}

	playURLGenerator func(challengeID uuid.UUID) string

	questionsService interface {
		GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*questions.Question, error)
		CheckAnswer(ctx context.Context, id uuid.UUID) (bool, error)
	}

	Question struct {
		QuestionID     string         `json:"question_id"`
		QuestionText   string         `json:"question_text"`
		TimeForAnswer  int            `json:"time_for_answer"`
		TotalQuestions int            `json:"total_questions"`
		QuestionNumber int            `json:"question_number"`
		AnswerOptions  []AnswerOption `json:"answer_options"`
	}

	AnswerOption struct {
		AnswerID   string `json:"answer_id"`
		AnswerText string `json:"answer_text"`
	}
)

// DefaultPlayURLGenerator ...
func DefaultPlayURLGenerator(baseURL string) playURLGenerator {
	return func(challengeID uuid.UUID) string {
		return fmt.Sprintf("%s/%s/play", baseURL, challengeID.String())
	}
}

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(cr challengesRepository, fn playURLGenerator, qs questionsService) *Service {
	if cr == nil {
		log.Fatalln("challenges repository is not set")
	}

	if qs == nil {
		log.Fatalln("questions service client is not set")
	}

	return &Service{
		cr:        cr,
		playUrlFn: fn,
		qs:        qs,
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

// GetVerificationQuestionByEpisodeID ...
func (s *Service) GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID uuid.UUID) (interface{}, error) {
	challenge, err := s.cr.GetChallengeByEpisodeID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}
	q, err := s.qs.GetOneRandomQuestionByChallengeID(ctx, challenge.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	answers := make([]AnswerOption, 0, len(q.AnswerOptions))
	for _, o := range q.AnswerOptions {
		answers = append(answers, AnswerOption{
			AnswerID:   o.ID.String(),
			AnswerText: o.Option,
		})
	}

	return Question{
		QuestionID:    q.ID.String(),
		QuestionText:  q.Question,
		TimeForAnswer: int(challenge.TimePerQuestion.Int32),
		AnswerOptions: answers,
	}, nil
}

// CheckVerificationQuestionAnswer ...
func (s *Service) CheckVerificationQuestionAnswer(ctx context.Context, qid, aid uuid.UUID) (interface{}, error) {
	return s.qs.CheckAnswer(ctx, aid) // FIXME: check question + answer, not only answer
}

// VerifyUserAccessToEpisode ...
func (s *Service) VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (interface{}, error) {
	// return s.qs.CheckAnswer(ctx, aid) // FIXME: check question + answer, not only answer
	return false, nil
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
		ShowID:             c.ShowID,
		Title:              c.Title,
		Description:        c.Description.String,
		PrizePool:          fmt.Sprintf("%.2f SAO", c.PrizePool),
		PrizePoolAmount:    c.PrizePool,
		Players:            c.PlayersToStart,
		TimePerQuestion:    fmt.Sprintf("%d sec", c.TimePerQuestion.Int32),
		TimePerQuestionSec: int64(c.TimePerQuestion.Int32),
		Play:               playUrlFn(c.ID),
		EpisodeID:          c.EpisodeID,
		Kind:               c.Kind,
	}
}

// AddChallenge ..
func (s *Service) AddChallenge(ctx context.Context, ch Challenge) (Challenge, error) {
	challenge, err := s.cr.AddChallenge(ctx, repository.AddChallengeParams{
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Description,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      ch.PrizePoolAmount,
		PlayersToStart: ch.Players,
		TimePerQuestion: sql.NullInt32{
			Int32: int32(ch.TimePerQuestionSec),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		EpisodeID: ch.EpisodeID,
		Kind:      ch.Kind,
	})
	if err != nil {
		return Challenge{}, fmt.Errorf("could not add challenge with title=%s: %w", ch.Title, err)
	}

	return castToChallenge(challenge, s.playUrlFn), nil
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
	if err := s.cr.UpdateChallenge(ctx, repository.UpdateChallengeParams{
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Title,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      ch.PrizePoolAmount,
		PlayersToStart: ch.Players,
		TimePerQuestion: sql.NullInt32{
			Int32: int32(ch.TimePerQuestionSec),
			Valid: false,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		EpisodeID: ch.EpisodeID,
		Kind:      ch.Kind,
		ID:        ch.ID,
	}); err != nil {
		return fmt.Errorf("could not update challenge with id=%s:%w", ch.ID, err)
	}

	return nil
}
