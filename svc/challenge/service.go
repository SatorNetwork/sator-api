package challenge

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
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

	challengesRepository interface {
		AddChallenge(ctx context.Context, arg repository.AddChallengeParams) (repository.Challenge, error)
		GetChallenges(ctx context.Context, arg repository.GetChallengesParams) ([]repository.Challenge, error)
		GetChallengeByEpisodeID(ctx context.Context, episodeID uuid.UUID) (repository.Challenge, error)
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository.Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, arg repository.UpdateChallengeParams) error

		AddQuestion(ctx context.Context, arg repository.AddQuestionParams) (repository.Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (repository.Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]repository.Question, error)
		UpdateQuestion(ctx context.Context, arg repository.UpdateQuestionParams) error

		AddQuestionOption(ctx context.Context, arg repository.AddQuestionOptionParams) (repository.AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, arg repository.DeleteAnswerByIDParams) error
		GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]repository.AnswerOption, error)
		GetAnswersByIDs(ctx context.Context, questionIds []uuid.UUID) ([]repository.AnswerOption, error)
		CheckAnswer(ctx context.Context, arg repository.CheckAnswerParams) (sql.NullBool, error)
		UpdateAnswer(ctx context.Context, arg repository.UpdateAnswerParams) error
	}

	playURLGenerator func(challengeID uuid.UUID) string

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

	// Question struct
	Question struct {
		ID             uuid.UUID      `json:"id"`
		ChallengeID    uuid.UUID      `json:"challenge_id"`
		Question       string         `json:"question"`
		TimeForAnswer  int            `json:"time_for_answer"`
		TotalQuestions int            `json:"total_questions"` // TODO: Do we need this field?
		Order          int32          `json:"order"`
		AnswerOptions  []AnswerOption `json:"answer_options"`
	}

	AnswerOption struct {
		ID         uuid.UUID `json:"id"`
		QuestionID uuid.UUID `json:"question_id"`
		Option     string    `json:"option"`
		IsCorrect  bool      `json:"is_correct"`
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

// GetVerificationQuestionByEpisodeID ...
// TODO: THIS METHOD
func (s *Service) GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID uuid.UUID) (interface{}, error) {
	// TODO: How many attempt

	challenge, err := s.cr.GetChallengeByEpisodeID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	q, err := s.GetOneRandomQuestionByChallengeID(ctx, challenge.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	answers := make([]AnswerOption, 0, len(q.AnswerOptions))
	for _, o := range q.AnswerOptions {
		answers = append(answers, AnswerOption{
			ID:         o.ID,
			QuestionID: o.QuestionID,
			Option:     o.Option,
			IsCorrect:  o.IsCorrect,
		})
	}

	return Question{
		ID:            q.ID,
		ChallengeID:   q.ChallengeID,
		Question:      q.Question,
		TimeForAnswer: int(challenge.TimePerQuestion.Int32),
		Order:         q.Order,
		AnswerOptions: answers,
	}, nil
}

// CheckVerificationQuestionAnswer ...
// TODO: THIS METHOD
func (s *Service) CheckVerificationQuestionAnswer(ctx context.Context, qid, aid uuid.UUID) (interface{}, error) {

	// TODO: Check how many attempts made
	// attempt, err := GetAttemptAmount
	// if attempt >= 2{return nil, errors.New("user has no more attempts to pass verification question")
	isValid, err := s.CheckAnswer(ctx, aid, qid)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}
	if isValid == false {
		// TODO: store failed attempt
	}

	return isValid, nil
}

// VerifyUserAccessToEpisode ...
// TODO: THIS METHOD
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

// GetQuestionByID returns question by id
func (s *Service) GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error) {
	question, err := s.cr.GetQuestionByID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Question{}, fmt.Errorf("could not get question: %w", err)
		}
		return Question{}, fmt.Errorf("could not found question with id=%s: %w", id.String(), err)
	}

	answers, err := s.cr.GetAnswersByQuestionID(ctx, question.ID)
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
	questions, err := s.cr.GetQuestionsByChallengeID(ctx, id)
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

	answers, err := s.cr.GetAnswersByIDs(ctx, idsSlice)
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

// GetOneRandomQuestionByChallengeID returns one random question by challenge id
func (s *Service) GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*Question, error) {
	questions, err := s.cr.GetQuestionsByChallengeID(ctx, id)
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

	answers, err := s.cr.GetAnswersByIDs(ctx, idsSlice)
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

	return &qlist[rand.Intn(len(qlist)-1)], nil
}

// CheckAnswer checks answer
func (s *Service) CheckAnswer(ctx context.Context, aid, qid uuid.UUID) (bool, error) {
	answers, err := s.cr.CheckAnswer(ctx, repository.CheckAnswerParams{
		ID:         aid,
		QuestionID: qid,
	})
	if err != nil {
		if !db.IsNotFoundError(err) {
			return false, fmt.Errorf("could not validate answer: %w", err)
		}
		return false, fmt.Errorf("could not found answer option with id %s: %w", aid, err)
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
	question, err := s.cr.AddQuestion(ctx, repository.AddQuestionParams{
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
	answer, err := s.cr.AddQuestionOption(ctx, repository.AddQuestionOptionParams{
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
	if err := s.cr.DeleteQuestionByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete question with id=%s:%w", id, err)
	}

	return nil
}

// DeleteAnswerByID ...
func (s *Service) DeleteAnswerByID(ctx context.Context, id, questionID uuid.UUID) error {
	if err := s.cr.DeleteAnswerByID(ctx, repository.DeleteAnswerByIDParams{
		ID:         id,
		QuestionID: questionID,
	}); err != nil {
		return fmt.Errorf("could not delete answer with id=%s:%w", id, err)
	}

	return nil
}

// UpdateQuestion ..
func (s *Service) UpdateQuestion(ctx context.Context, qw Question) error {
	if err := s.cr.UpdateQuestion(ctx, repository.UpdateQuestionParams{
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
	if err := s.cr.UpdateAnswer(ctx, repository.UpdateAnswerParams{
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
