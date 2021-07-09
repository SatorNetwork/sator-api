package questions

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/questions/repository"
	"github.com/google/uuid"
)

type (
	//Service struct
	Service struct {
		qr questionsRepository
	}

	questionsRepository interface {
		AddQuestion(ctx context.Context, arg repository.AddQuestionParams) (repository.Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (repository.Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]repository.Question, error)
		UpdateQuestion(ctx context.Context, arg repository.UpdateQuestionParams) error

		AddQuestionOption(ctx context.Context, arg repository.AddQuestionOptionParams) (repository.AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, id uuid.UUID) error
		GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]repository.AnswerOption, error)
		GetAnswersByIDs(ctx context.Context, questionIds []uuid.UUID) ([]repository.AnswerOption, error)
		CheckAnswer(ctx context.Context, id uuid.UUID) (sql.NullBool, error)
		UpdateAnswer(ctx context.Context, arg repository.UpdateAnswerParams) error
	}

	// Question struct
	Question struct {
		ID            uuid.UUID      `json:"id"`
		ChallengeID   uuid.UUID      `json:"challenge_id"`
		Question      string         `json:"question"`
		Order         int32          `json:"order"`
		AnswerOptions []AnswerOption `json:"question_options"`
	}

	// AnswerOption struct
	AnswerOption struct {
		ID         uuid.UUID `json:"id"`
		QuestionID uuid.UUID `json:"question_id"`
		Option     string    `json:"option"`
		IsCorrect  bool      `json:"is_correct"`
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(qr questionsRepository) *Service {
	if qr == nil {
		log.Fatalln("questions repository is not set")
	}

	return &Service{qr: qr}
}

// GetQuestionByID returns question by id
func (s *Service) GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error) {
	question, err := s.qr.GetQuestionByID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Question{}, fmt.Errorf("could not get question: %w", err)
		}
		return Question{}, fmt.Errorf("could not found question with i=%s: %w", id.String(), err)
	}

	answers, err := s.qr.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Question{}, fmt.Errorf("could not get answer options for question with id=%s: %w", id.String(), err)
		}

		return Question{}, fmt.Errorf("could not found any answer options for question with id=%s: %w", id.String(), err)
	}

	return castToQuestionWithAnswers(question, answers), nil
}

// GetQuestionsByChallengeID returns questions by challenge id
func (s *Service) GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
		}
		return nil, fmt.Errorf("could not found any questions with challenge id %s: %w", id.String(), err)
	}

	idsSlice := make([]uuid.UUID, 0, len(questions))
	for _, v := range questions {
		idsSlice = append(idsSlice, v.ID)
	}

	answers, err := s.qr.GetAnswersByIDs(ctx, idsSlice)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get answers: %w", err)
		}
		return nil, fmt.Errorf("could not found answers with ids %s: %w", id.String(), err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	answMap := make(map[string][]AnswerOption)
	for _, v := range answers {
		if _, ok := answMap[v.QuestionID.String()]; ok {
			answMap[v.QuestionID.String()] = append(answMap[v.QuestionID.String()], AnswerOption{
				ID:         v.ID,
				QuestionID: v.QuestionID,
				Option:     v.AnswerOption,
				IsCorrect:  v.IsCorrect.Bool,
			})
		} else {
			answMap[v.QuestionID.String()] = []AnswerOption{
				{
					ID:         v.ID,
					QuestionID: v.QuestionID,
					Option:     v.AnswerOption,
					IsCorrect:  v.IsCorrect.Bool,
				},
			}
		}
	}

	qlist := make([]Question, 0, len(questions))
	for _, v := range questions {
		if opt, ok := answMap[v.ID.String()]; ok {
			qlist = append(qlist, Question{
				ID:            v.ID,
				ChallengeID:   v.ChallengeID,
				Question:      v.Question,
				Order:         v.QuestionOrder,
				AnswerOptions: opt,
			})
		}
	}

	return qlist, nil
}

// CheckAnswer checks answer
func (s *Service) CheckAnswer(ctx context.Context, id uuid.UUID) (bool, error) {
	answers, err := s.qr.CheckAnswer(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return false, fmt.Errorf("could not validate answer: %w", err)
		}
		return false, fmt.Errorf("could not found question with id %s: %w", id, err)
	}

	return answers.Bool, nil
}

func castToQuestionWithAnswers(q repository.Question, a []repository.AnswerOption) Question {
	options := make([]AnswerOption, 0, len(a))
	for _, ao := range a {
		options = append(options, AnswerOption{
			ID:         ao.ID,
			QuestionID: ao.QuestionID,
			Option:     ao.AnswerOption,
			IsCorrect:  ao.IsCorrect.Bool,
		})
	}

	return Question{
		ID:            q.ID,
		ChallengeID:   q.ChallengeID,
		Question:      q.Question,
		Order:         q.QuestionOrder,
		AnswerOptions: options,
	}
}

// AddQuestion ..
func (s *Service) AddQuestion(ctx context.Context, qw Question) (Question, error) {
	question, err := s.qr.AddQuestion(ctx, repository.AddQuestionParams{
		ChallengeID:   qw.ChallengeID,
		Question:      qw.Question,
		QuestionOrder: qw.Order,
	})
	if err != nil {
		return Question{}, fmt.Errorf("could not add question %s: %v", qw.Question, err)
	}

	return castToQuestion(question), nil
}

func castToQuestion(q repository.Question) Question {
	return Question{
		ID:          q.ID,
		ChallengeID: q.ChallengeID,
		Question:    q.Question,
		Order:       q.QuestionOrder,
	}
}

// AddQuestionOption ..
func (s *Service) AddQuestionOption(ctx context.Context, ao AnswerOption) (AnswerOption, error) {
	answer, err := s.qr.AddQuestionOption(ctx, repository.AddQuestionOptionParams{
		QuestionID:   ao.QuestionID,
		AnswerOption: ao.Option,
		IsCorrect: sql.NullBool{
			Bool:  ao.IsCorrect,
			Valid: true,
		},
	})
	if err != nil {
		return AnswerOption{}, fmt.Errorf("could not add answer %s: %v", ao.Option, err)
	}

	return castToAnswerOption(answer), nil
}

func castToAnswerOption(ao repository.AnswerOption) AnswerOption {
	return AnswerOption{
		ID:         ao.ID,
		QuestionID: ao.QuestionID,
		Option:     ao.AnswerOption,
		IsCorrect:  ao.IsCorrect.Bool,
	}
}

// DeleteQuestionByID ...
func (s *Service) DeleteQuestionByID(ctx context.Context, id uuid.UUID) error {
	if err := s.qr.DeleteQuestionByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete question with id=%s:%w", id, err)
	}

	return nil
}

// DeleteAnswerByID ...
func (s *Service) DeleteAnswerByID(ctx context.Context, id uuid.UUID) error {
	if err := s.qr.DeleteQuestionByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete answer with id=%s:%w", id, err)
	}

	return nil
}

// UpdateQuestion ..
func (s *Service) UpdateQuestion(ctx context.Context, qw Question) error {
	if err := s.qr.UpdateQuestion(ctx, repository.UpdateQuestionParams{
		ID:            qw.ID,
		ChallengeID:   qw.ChallengeID,
		Question:      qw.Question,
		QuestionOrder: qw.Order,
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not update question with id=%s:%w", qw.ID, err)
	}

	return nil
}

// UpdateAnswer ..
func (s *Service) UpdateAnswer(ctx context.Context, ao AnswerOption) error {
	if err := s.qr.UpdateAnswer(ctx, repository.UpdateAnswerParams{
		ID:           ao.ID,
		QuestionID:   ao.QuestionID,
		AnswerOption: ao.Option,
		IsCorrect: sql.NullBool{
			Bool:  ao.IsCorrect,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not update answer with id=%s:%w", ao.ID, err)
	}

	return nil
}
