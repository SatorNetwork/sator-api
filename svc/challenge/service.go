package challenge

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/lib/db"
	"github.com/SatorNetwork/sator-api/svc/challenge/repository"
	showRepository "github.com/SatorNetwork/sator-api/svc/shows/repository"

	"github.com/google/uuid"
)

const (
	defaultMaxWinners        = 1
	defaultQuestionsPerGame  = 5
	defaultMinCorrectAnswers = 1

	defaultPercentForQuiz = 5
	defaultMinimumReward  = 1
)

type (
	// Service struct
	Service struct {
		cr                   challengesRepository
		playUrlFn            playURLGenerator
		attemptsNumber       int64
		activatedRealmPeriod time.Duration
		chargeForUnlockFn    chargeForUnlockFunc
		showRepo             showsRepository
	}

	chargeForUnlockFunc func(ctx context.Context, uid uuid.UUID, amount float64, info string) error

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	// FIXME: remove it when the app will be fixed
	showsRepository interface {
		GetEpisodeByID(ctx context.Context, id uuid.UUID) (showRepository.GetEpisodeByIDRow, error)
		GetEpisodeIDByVerificationChallengeID(ctx context.Context, verificationChallengeID uuid.NullUUID) (uuid.UUID, error)
		GetEpisodeIDByQuizChallengeID(ctx context.Context, quizChallengeID uuid.NullUUID) (uuid.UUID, error)
	}

	challengesRepository interface {
		AddChallenge(ctx context.Context, arg repository.AddChallengeParams) (repository.Challenge, error)
		GetChallenges(ctx context.Context, arg repository.GetChallengesParams) ([]repository.Challenge, error)
		GetChallengeByEpisodeID(ctx context.Context, episodeID uuid.NullUUID) (repository.Challenge, error)
		GetChallengeByID(ctx context.Context, id uuid.UUID) (repository.Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, arg repository.UpdateChallengeParams) error

		// Questions
		AddQuestion(ctx context.Context, arg repository.AddQuestionParams) (repository.Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (repository.Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]repository.Question, error)
		GetQuestionsByChallengeIDWithExceptions(ctx context.Context, arg repository.GetQuestionsByChallengeIDWithExceptionsParams) ([]repository.Question, error)
		UpdateQuestion(ctx context.Context, arg repository.UpdateQuestionParams) error

		// Answers
		AddQuestionOption(ctx context.Context, arg repository.AddQuestionOptionParams) (repository.AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, arg repository.DeleteAnswerByIDParams) error
		DeleteAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) error
		GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]repository.AnswerOption, error)
		GetAnswersByIDs(ctx context.Context, questionIds []uuid.UUID) ([]repository.AnswerOption, error)
		CheckAnswer(ctx context.Context, arg repository.CheckAnswerParams) (sql.NullBool, error)
		UpdateAnswer(ctx context.Context, arg repository.UpdateAnswerParams) error

		// Episode Access
		AddEpisodeAccessData(ctx context.Context, arg repository.AddEpisodeAccessDataParams) (repository.EpisodeAccess, error)
		DeleteEpisodeAccessData(ctx context.Context, arg repository.DeleteEpisodeAccessDataParams) error
		GetEpisodeAccessData(ctx context.Context, arg repository.GetEpisodeAccessDataParams) (repository.EpisodeAccess, error)
		UpdateEpisodeAccessData(ctx context.Context, arg repository.UpdateEpisodeAccessDataParams) error
		DoesUserHaveAccessToEpisode(ctx context.Context, arg repository.DoesUserHaveAccessToEpisodeParams) (bool, error)
		NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error)
		ListIDsAvailableUserEpisodes(ctx context.Context, arg repository.ListIDsAvailableUserEpisodesParams) ([]uuid.UUID, error)

		// Verification Question Attempts
		AddAttempt(ctx context.Context, arg repository.AddAttemptParams) (repository.Attempt, error)
		GetEpisodeIDByQuestionID(ctx context.Context, arg repository.GetEpisodeIDByQuestionIDParams) (uuid.UUID, error)
		CountAttempts(ctx context.Context, arg repository.CountAttemptsParams) (int64, error)
		GetAskedQuestionsByEpisodeID(ctx context.Context, arg repository.GetAskedQuestionsByEpisodeIDParams) ([]uuid.UUID, error)
		UpdateAttempt(ctx context.Context, arg repository.UpdateAttemptParams) error

		// Challenge Attempts
		AddChallengeAttempt(ctx context.Context, arg repository.AddChallengeAttemptParams) (repository.PassedChallengesDatum, error)
		StoreChallengeReceivedRewardAmount(ctx context.Context, arg repository.StoreChallengeReceivedRewardAmountParams) error
		CountPassedChallengeAttempts(ctx context.Context, arg repository.CountPassedChallengeAttemptsParams) (int64, error)
		GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error)
		GetChallengeReceivedRewardAmountByUserID(ctx context.Context, arg repository.GetChallengeReceivedRewardAmountByUserIDParams) (float64, error)
	}

	playURLGenerator func(challengeID uuid.UUID) string

	// Challenge struct
	// Fields were rearranged to optimize memory usage.
	Challenge struct {
		ID                     uuid.UUID  `json:"id"`
		ShowID                 uuid.UUID  `json:"show_id"`
		Title                  string     `json:"title"`
		Description            string     `json:"description"`
		PrizePool              string     `json:"prize_pool"`
		PrizePoolAmount        float64    `json:"prize_pool_amount"`
		Players                int32      `json:"players"`
		Winners                string     `json:"winners"`
		TimePerQuestion        string     `json:"time_per_question"`
		TimePerQuestionSec     int32      `json:"time_per_question_sec"`
		Play                   string     `json:"play"`
		EpisodeID              *uuid.UUID `json:"episode_id"`
		Kind                   int32      `json:"kind"`
		UserMaxAttempts        int32      `json:"user_max_attempts"`
		AttemptsLeft           int32      `json:"attempts_left"`
		ReceivedReward         float64    `json:"received_reward"`
		ReceivedRewardStr      string     `json:"received_reward_str"`
		MaxWinners             int32      `json:"max_winners"`
		QuestionsPerGame       int32      `json:"questions_per_game"`
		MinCorrectAnswers      int32      `json:"min_correct_answers"`
		IsRealmActivated       bool       `json:"is_realm_activated"`
		RegisteredPlayers      int        `json:"registered_players"`
		RegisteredPlayersInDB  int        `json:"registered_players_in_db"`
		PercentForQuiz         float64    `json:"percent_for_quiz"`
		MinimumReward          float64    `json:"minimum_reward"`
		CurrentPrizePool       string     `json:"current_prize_pool"`
		CurrentPrizePoolAmount float64    `json:"current_prize_pool_amount"`
	}

	RawChallenge struct {
		ID                 uuid.UUID  `json:"id"`
		ShowID             uuid.UUID  `json:"show_id"`
		Title              string     `json:"title"`
		Description        string     `json:"description"`
		PrizePool          string     `json:"prize_pool"`
		PrizePoolAmount    float64    `json:"prize_pool_amount"`
		PlayersToStart     int32      `json:"players_to_start"`
		TimePerQuestion    string     `json:"time_per_question"`
		TimePerQuestionSec int32      `json:"time_per_question_sec"`
		EpisodeID          *uuid.UUID `json:"episode_id"`
		Kind               int32      `json:"kind"`
		UserMaxAttempts    int32      `json:"user_max_attempts"`
		MaxWinners         int32      `json:"max_winners"`
		QuestionsPerGame   int32      `json:"questions_per_game"`
		MinCorrectAnswers  int32      `json:"min_correct_answers"`
		PercentForQuiz     float64    `json:"percent_for_quiz"`
		MinimumReward      float64    `json:"minimum_reward"`
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

	// AnswerOption struct
	AnswerOption struct {
		ID         uuid.UUID `json:"id"`
		QuestionID uuid.UUID `json:"question_id"`
		Option     string    `json:"option"`
		IsCorrect  bool      `json:"is_correct"`
	}

	QuizQuestion struct {
		QuestionID     string             `json:"question_id"`
		QuestionText   string             `json:"question_text"`
		TimeForAnswer  int                `json:"time_for_answer"`
		TotalQuestions int                `json:"total_questions"`
		QuestionNumber int                `json:"question_number"`
		AnswerOptions  []QuizAnswerOption `json:"answer_options"`
	}

	QuizAnswerOption struct {
		AnswerID   string `json:"answer_id"`
		AnswerText string `json:"answer_text"`
	}

	EpisodeAccess struct {
		EpisodeID       *uuid.UUID `json:"episode_id"`
		Result          bool       `json:"result"`
		ActivatedAt     string     `json:"activated_at,omitempty"`
		ActivatedBefore string     `json:"activated_before,omitempty"`
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
func NewService(cr challengesRepository, showRepo showsRepository, fn playURLGenerator, opt ...ServiceOption) *Service {
	if cr == nil {
		log.Fatalln("challenges repository is not set")
	}

	s := &Service{
		cr:                   cr,
		showRepo:             showRepo,
		playUrlFn:            fn,
		attemptsNumber:       2,
		activatedRealmPeriod: time.Hour * 24,
	}

	for _, o := range opt {
		o(s)
	}

	return s
}

// GetByID ...
func (s *Service) GetByID(ctx context.Context, challengeID, userID uuid.UUID) (Challenge, error) {
	challenge, err := s.cr.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return Challenge{}, fmt.Errorf("could not get challenge by challengeID=%s: %w", challengeID, err)
	}

	var attemptsLeft int32

	receivedReward, err := s.cr.GetChallengeReceivedRewardAmountByUserID(ctx, repository.GetChallengeReceivedRewardAmountByUserIDParams{
		UserID:      userID,
		ChallengeID: challengeID,
	})
	if err != nil && !db.IsNotFoundError(err) {
		return Challenge{}, fmt.Errorf("could not get received reward amount: %w", err)
	}

	if receivedReward == 0 {
		attempts, err := s.cr.CountPassedChallengeAttempts(ctx, repository.CountPassedChallengeAttemptsParams{
			UserID:      userID,
			ChallengeID: challengeID,
		})
		if err != nil {
			return Challenge{}, fmt.Errorf("could not get passed challenge attempts: %w", err)
		}
		attemptsLeft = challenge.UserMaxAttempts - int32(attempts)
		if attemptsLeft < 0 {
			attemptsLeft = 0
		}
	}

	// Check if a user has access to the challenge
	epID, _ := s.showRepo.GetEpisodeIDByQuizChallengeID(ctx, uuid.NullUUID{UUID: challengeID, Valid: true})
	res, _ := s.VerifyUserAccessToEpisode(ctx, userID, epID)

	return castToChallenge(challenge, s.playUrlFn, attemptsLeft, receivedReward, &epID, res.Result), nil
}

func (s *Service) GetRawChallengeByID(ctx context.Context, challengeID uuid.UUID) (RawChallenge, error) {
	challenge, err := s.cr.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return RawChallenge{}, fmt.Errorf("could not get challenge by challengeID=%s: %w", challengeID, err)
	}

	return castToRawChallenge(challenge), nil
}

// GetChallengeByID ...
func (s *Service) GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (interface{}, error) {
	return s.GetByID(ctx, challengeID, userID)
}

// GetVerificationQuestionByEpisodeID ...
func (s *Service) GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID, userID uuid.UUID) (interface{}, error) {
	ep, err := s.showRepo.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not found episode with id=%s: %w", episodeID, err)
	}
	challenge, err := s.cr.GetChallengeByID(ctx, ep.VerificationChallengeID.UUID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	numberAttempts, err := s.cr.CountAttempts(ctx, repository.CountAttemptsParams{
		UserID:    userID,
		EpisodeID: episodeID,
		CreatedAt: sql.NullTime{
			Time:  time.Now().Add(-s.activatedRealmPeriod),
			Valid: true,
		},
	})
	if err == nil {
		if numberAttempts >= int64(challenge.UserMaxAttempts) {
			return nil, fmt.Errorf("user has no more attempts to pass verification question")
		}
	}

	askedQuestions, _ := s.cr.GetAskedQuestionsByEpisodeID(ctx, repository.GetAskedQuestionsByEpisodeIDParams{
		UserID:    userID,
		EpisodeID: episodeID,
	})
	q, err := s.GetOneRandomQuestionByChallengeID(ctx, challenge.ID, askedQuestions...)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge by id: %w", err)
	}

	// store attempt anyway
	if _, err := s.cr.AddAttempt(ctx, repository.AddAttemptParams{
		UserID:     userID,
		EpisodeID:  episodeID,
		QuestionID: q.ID,
	}); err != nil {
		return nil, fmt.Errorf("could not add attempt data: %w", err)
	}

	answers := make([]QuizAnswerOption, 0, len(q.AnswerOptions))
	for _, o := range q.AnswerOptions {
		answers = append(answers, QuizAnswerOption{
			AnswerID:   o.ID.String(),
			AnswerText: o.Option,
		})
	}

	return QuizQuestion{
		QuestionID:    q.ID.String(),
		QuestionText:  q.Question,
		TimeForAnswer: int(challenge.TimePerQuestion.Int32),
		AnswerOptions: answers,
	}, nil
}

// CheckVerificationQuestionAnswer ...
func (s *Service) CheckVerificationQuestionAnswer(ctx context.Context, questionID, answerID, userID uuid.UUID) (interface{}, error) {
	question, err := s.cr.GetQuestionByID(ctx, questionID)
	if err != nil {
		return nil, fmt.Errorf("could not get question by id: %w", err)
	}

	// challenge, err := s.cr.GetChallengeByID(ctx, question.ChallengeID)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get challenge by id: %w", err)
	// }

	isValid, err := s.CheckAnswer(ctx, answerID, questionID)
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}

	if err := s.cr.UpdateAttempt(ctx, repository.UpdateAttemptParams{
		AnswerID:   uuid.NullUUID{UUID: answerID, Valid: answerID != uuid.Nil},
		Valid:      sql.NullBool{Bool: isValid, Valid: true},
		UserID:     userID,
		QuestionID: questionID,
	}); err != nil {
		return nil, err
	}

	if !isValid {
		return false, nil
	}

	epID, err := s.showRepo.GetEpisodeIDByVerificationChallengeID(ctx, uuid.NullUUID{UUID: question.ChallengeID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("could not get episode id: %w", err)
	}

	if _, err := s.cr.AddEpisodeAccessData(ctx, repository.AddEpisodeAccessDataParams{
		EpisodeID: epID,
		UserID:    userID,
		ActivatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ActivatedBefore: sql.NullTime{
			Time:  time.Now().Add(s.activatedRealmPeriod),
			Valid: true,
		},
	}); err != nil {
		return false, fmt.Errorf("could not store episode access data: %w", err)
	}

	return true, nil
}

// VerifyUserAccessToEpisode ...
func (s *Service) VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (EpisodeAccess, error) {
	data, err := s.cr.GetEpisodeAccessData(ctx, repository.GetEpisodeAccessDataParams{
		EpisodeID: eid,
		UserID:    uid,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return EpisodeAccess{Result: false}, nil
		}
		return EpisodeAccess{Result: false}, fmt.Errorf("could not get episode access data: %w", err)
	}

	if !data.ActivatedAt.Valid || !data.ActivatedBefore.Valid || data.ActivatedBefore.Time.Before(time.Now()) {
		return EpisodeAccess{Result: false}, nil
	}

	return EpisodeAccess{
		Result:          data.ActivatedBefore.Time.After(time.Now()),
		ActivatedAt:     data.ActivatedAt.Time.Format(time.RFC3339),
		ActivatedBefore: data.ActivatedBefore.Time.Format(time.RFC3339),
	}, nil
}

// GetChallengesByShowID ...
func (s *Service) GetChallengesByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error) {
	list, err := s.cr.GetChallenges(ctx, repository.GetChallengesParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get challenge list by show id: %w", err)
	}
	var attemptsLeft int32
	// Cast repository.Challenge into challenge.Challenge struct
	result := make([]Challenge, 0, len(list))
	for _, v := range list {
		receivedReward, err := s.cr.GetChallengeReceivedRewardAmountByUserID(ctx, repository.GetChallengeReceivedRewardAmountByUserIDParams{
			UserID:      userID,
			ChallengeID: v.ID,
		})
		if err != nil && !db.IsNotFoundError(err) {
			return Challenge{}, fmt.Errorf("could not get received reward amount: %w", err)
		}

		if receivedReward == 0 {
			attempts, err := s.cr.CountPassedChallengeAttempts(ctx, repository.CountPassedChallengeAttemptsParams{
				UserID:      userID,
				ChallengeID: v.ID,
			})
			if err != nil {
				return Challenge{}, fmt.Errorf("could not get passed challenge attempts: %w", err)
			}
			attemptsLeft = v.UserMaxAttempts - int32(attempts)
			if attemptsLeft < 0 {
				attemptsLeft = 0
			}
		}
		result = append(result, castToChallenge(v, s.playUrlFn, attemptsLeft, receivedReward, nil, false))
	}

	return result, nil
}

func castToChallenge(c repository.Challenge, playUrlFn playURLGenerator, attemptsLeft int32, receivedReward float64, epID *uuid.UUID, isActivated bool) Challenge {
	ch := Challenge{
		ID:                 c.ID,
		ShowID:             c.ShowID,
		EpisodeID:          epID,
		Title:              c.Title,
		Description:        c.Description.String,
		PrizePool:          fmt.Sprintf("%.2f SAO", c.PrizePool),
		PrizePoolAmount:    c.PrizePool,
		Players:            c.PlayersToStart,
		TimePerQuestion:    fmt.Sprintf("%d sec", c.TimePerQuestion.Int32),
		TimePerQuestionSec: c.TimePerQuestion.Int32,
		Play:               playUrlFn(c.ID),
		Kind:               c.Kind,
		UserMaxAttempts:    c.UserMaxAttempts,
		AttemptsLeft:       attemptsLeft,
		ReceivedReward:     receivedReward,
		ReceivedRewardStr:  fmt.Sprintf("%.2f SAO", receivedReward),
		MaxWinners:         c.MaxWinners.Int32,
		QuestionsPerGame:   c.QuestionsPerGame,
		MinCorrectAnswers:  c.MinCorrectAnswers,
		IsRealmActivated:   isActivated,
		PercentForQuiz:     c.PercentForQuiz,
		MinimumReward:      c.MinimumReward,
	}

	if ch.MaxWinners == 0 {
		ch.MaxWinners = defaultMaxWinners
	}

	if c.EpisodeID.Valid && c.EpisodeID.UUID != uuid.Nil {
		ch.EpisodeID = &c.EpisodeID.UUID
	}

	return ch
}

func castToRawChallenge(c repository.Challenge) RawChallenge {
	rawChallenge := RawChallenge{
		ID:                 c.ID,
		ShowID:             c.ShowID,
		Title:              c.Title,
		Description:        c.Description.String,
		PrizePool:          fmt.Sprintf("%.2f SAO", c.PrizePool),
		PrizePoolAmount:    c.PrizePool,
		PlayersToStart:     c.PlayersToStart,
		TimePerQuestion:    fmt.Sprintf("%d sec", c.TimePerQuestion.Int32),
		TimePerQuestionSec: c.TimePerQuestion.Int32,
		Kind:               c.Kind,
		UserMaxAttempts:    c.UserMaxAttempts,
		MaxWinners:         c.MaxWinners.Int32,
		QuestionsPerGame:   c.QuestionsPerGame,
		MinCorrectAnswers:  c.MinCorrectAnswers,
		PercentForQuiz:     c.PercentForQuiz,
		MinimumReward:      c.MinimumReward,
	}

	if c.EpisodeID.Valid && c.EpisodeID.UUID != uuid.Nil {
		rawChallenge.EpisodeID = &c.EpisodeID.UUID
	}

	return rawChallenge
}

// AddChallenge ..
func (s *Service) AddChallenge(ctx context.Context, ch Challenge) (Challenge, error) {
	if ch.MinCorrectAnswers > ch.QuestionsPerGame {
		return Challenge{}, fmt.Errorf("min correct answers should be less or equal to questings per game")
	}

	if ch.MaxWinners == 0 {
		ch.MaxWinners = defaultMaxWinners
	}
	if ch.QuestionsPerGame == 0 {
		ch.QuestionsPerGame = defaultQuestionsPerGame
	}
	if ch.MinCorrectAnswers == 0 {
		ch.MinCorrectAnswers = defaultMinCorrectAnswers
	}
	if ch.PercentForQuiz == 0 {
		ch.PercentForQuiz = defaultPercentForQuiz
	}
	if ch.MinimumReward == 0 {
		ch.MinimumReward = defaultMinimumReward
	}

	params := repository.AddChallengeParams{
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Description,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      ch.PrizePoolAmount,
		PlayersToStart: ch.Players,
		TimePerQuestion: sql.NullInt32{
			Int32: ch.TimePerQuestionSec,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		Kind:            ch.Kind,
		UserMaxAttempts: ch.UserMaxAttempts,
		MaxWinners: sql.NullInt32{
			Int32: ch.MaxWinners,
			Valid: true,
		},
		QuestionsPerGame:  ch.QuestionsPerGame,
		MinCorrectAnswers: ch.MinCorrectAnswers,
		PercentForQuiz:    ch.PercentForQuiz,
		MinimumReward:     ch.MinimumReward,
	}

	if ch.EpisodeID != nil && *ch.EpisodeID != uuid.Nil {
		params.EpisodeID = uuid.NullUUID{UUID: *ch.EpisodeID, Valid: true}
	}

	challenge, err := s.cr.AddChallenge(ctx, params)
	if err != nil {
		return Challenge{}, fmt.Errorf("could not add challenge with title=%s: %w", ch.Title, err)
	}

	return castToChallenge(challenge, s.playUrlFn, 0, 0, nil, false), nil
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
	if ch.MinCorrectAnswers > ch.QuestionsPerGame {
		return fmt.Errorf("min correct answers should be less or equal to questings per game")
	}

	if ch.MaxWinners == 0 {
		ch.MaxWinners = defaultMaxWinners
	}
	if ch.QuestionsPerGame == 0 {
		ch.QuestionsPerGame = defaultQuestionsPerGame
	}
	if ch.MinCorrectAnswers == 0 {
		ch.MinCorrectAnswers = defaultMinCorrectAnswers
	}
	if ch.PercentForQuiz == 0 {
		ch.PercentForQuiz = defaultPercentForQuiz
	}
	if ch.MinimumReward == 0 {
		ch.MinimumReward = defaultMinimumReward
	}

	params := repository.UpdateChallengeParams{
		ID:     ch.ID,
		ShowID: ch.ShowID,
		Title:  ch.Title,
		Description: sql.NullString{
			String: ch.Description,
			Valid:  len(ch.Description) > 0,
		},
		PrizePool:      ch.PrizePoolAmount,
		PlayersToStart: ch.Players,
		TimePerQuestion: sql.NullInt32{
			Int32: ch.TimePerQuestionSec,
			Valid: ch.TimePerQuestionSec > 0,
		},
		Kind:            ch.Kind,
		UserMaxAttempts: ch.UserMaxAttempts,
		MaxWinners: sql.NullInt32{
			Int32: ch.MaxWinners,
			Valid: true,
		},
		QuestionsPerGame:  ch.QuestionsPerGame,
		MinCorrectAnswers: ch.MinCorrectAnswers,
		PercentForQuiz:    ch.PercentForQuiz,
		MinimumReward:     ch.MinimumReward,
	}

	if ch.EpisodeID != nil && *ch.EpisodeID != uuid.Nil {
		params.EpisodeID = uuid.NullUUID{UUID: *ch.EpisodeID, Valid: true}
	}

	if err := s.cr.UpdateChallenge(ctx, params); err != nil {
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
func (s *Service) GetOneRandomQuestionByChallengeID(ctx context.Context, challengeID uuid.UUID, excludeIDs ...uuid.UUID) (*Question, error) {
	var questions []repository.Question
	var err error

	if len(excludeIDs) > 0 {
		questions, err = s.cr.GetQuestionsByChallengeIDWithExceptions(ctx, repository.GetQuestionsByChallengeIDWithExceptionsParams{
			ChallengeID: challengeID,
			QuestionIds: excludeIDs,
		})
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
			}

			return nil, fmt.Errorf("could not found any questions with challenge id %s: %w", challengeID.String(), err)
		}
	} else {
		questions, err = s.cr.GetQuestionsByChallengeID(ctx, challengeID)
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
			}

			return nil, fmt.Errorf("could not found any questions with challenge id %s: %w", challengeID.String(), err)
		}
	}

	var q repository.Question
	switch len(questions) {
	case 0:
		return nil, fmt.Errorf("no more verification question")
	case 1:
		q = questions[0]
	default:
		q = questions[rand.Intn(len(questions))]
	}

	answers, err := s.cr.GetAnswersByQuestionID(ctx, q.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get answer options: %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	result := castToQuestionWithAnswers(q, answers)

	return &result, nil
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

// DeleteAnswersByQuestionID ...
func (s *Service) DeleteAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) error {
	if err := s.cr.DeleteAnswersByQuestionID(ctx, questionID); err != nil {
		return fmt.Errorf("could not delete answer options by question id=%s: %w", questionID, err)
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

// UnlockEpisode ...
func (s *Service) UnlockEpisode(ctx context.Context, userID, episodeID uuid.UUID, unlockOption string) error {
	activateBefore := time.Now()
	var amount float64

	if data, err := s.cr.GetEpisodeAccessData(ctx, repository.GetEpisodeAccessDataParams{
		EpisodeID: episodeID,
		UserID:    userID,
	}); err == nil && data.ActivatedBefore.Valid && data.ActivatedBefore.Time.After(time.Now()) {
		activateBefore = data.ActivatedBefore.Time
	}

	switch unlockOption {
	case "unlock_opt_2h":
		activateBefore = activateBefore.Add(time.Hour * 2)
		amount = 2
	case "unlock_opt_24h":
		activateBefore = activateBefore.Add(time.Hour * 24)
		amount = 10
	case "unlock_opt_week":
		activateBefore = activateBefore.Add(time.Hour * 24 * 7)
		amount = 50
	}

	if s.chargeForUnlockFn != nil && amount > 0 {
		if err := s.chargeForUnlockFn(
			ctx, userID, amount,
			fmt.Sprintf("unlock episode realm: %s", episodeID.String()),
		); err != nil {
			return fmt.Errorf("could not unlock episode realm: %w", err)
		}
	}

	if _, err := s.cr.AddEpisodeAccessData(ctx, repository.AddEpisodeAccessDataParams{
		EpisodeID: episodeID,
		UserID:    userID,
		ActivatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ActivatedBefore: sql.NullTime{
			Time:  activateBefore,
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not store episode access data: %w", err)
	}

	return nil
}

// StoreChallengeAttempt ...
func (s *Service) StoreChallengeAttempt(ctx context.Context, challengeID, userID uuid.UUID) error {
	if _, err := s.cr.AddChallengeAttempt(ctx, repository.AddChallengeAttemptParams{
		UserID:      userID,
		ChallengeID: challengeID,
	}); err != nil {
		return fmt.Errorf("could not store passed challenge data: %w", err)
	}

	return nil
}

// StoreChallengeReceivedRewardAmount ...
func (s *Service) StoreChallengeReceivedRewardAmount(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error {
	if err := s.cr.StoreChallengeReceivedRewardAmount(ctx, repository.StoreChallengeReceivedRewardAmountParams{
		UserID:       userID,
		ChallengeID:  challengeID,
		RewardAmount: rewardAmount,
	}); err != nil {
		return fmt.Errorf("could not store passed challenge data: %w", err)
	}

	return nil
}

// GetChallengeReceivedRewardAmountByUserID ...
func (s *Service) GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error) {
	amount, err := s.cr.GetChallengeReceivedRewardAmountByUserID(ctx, repository.GetChallengeReceivedRewardAmountByUserIDParams{
		UserID:      userID,
		ChallengeID: challengeID,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("could not get challenge received reward amount by user=%v: %w", userID, err)
	}

	return amount, nil
}

// GetChallengeReceivedRewardAmount ...
func (s *Service) GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error) {
	amount, err := s.cr.GetChallengeReceivedRewardAmount(ctx, challengeID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("could not get challenge received reward amount: %w", err)
	}

	return amount, nil
}

// GetPassedChallengeAttempts ...
func (s *Service) GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error) {
	attemptsNumber, err := s.cr.CountPassedChallengeAttempts(ctx, repository.CountPassedChallengeAttemptsParams{
		UserID:      userID,
		ChallengeID: challengeID,
	})
	if err != nil {
		return 0, fmt.Errorf("could not get passed challenge attempts: %w", err)
	}

	return attemptsNumber, nil
}

// NumberUsersWhoHaveAccessToEpisode ...
func (s *Service) NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error) {
	number, err := s.cr.NumberUsersWhoHaveAccessToEpisode(ctx, episodeID)
	if err != nil {
		return 0, fmt.Errorf("could not number users who have access to episode: %w", err)
	}

	return number, nil
}

// ListIDsAvailableUserEpisodes ...
func (s *Service) ListIDsAvailableUserEpisodes(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	list, err := s.cr.ListIDsAvailableUserEpisodes(ctx, repository.ListIDsAvailableUserEpisodesParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get list user available episodes: %w", err)
	}

	return list, nil
}

// GetAttemptsLeftForVerificationQuestion ...
func (s *Service) GetAttemptsLeftForVerificationQuestion(ctx context.Context, episodeID, userID uuid.UUID) (int64, error) {
	ep, err := s.showRepo.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return 0, fmt.Errorf("could not found episode with id=%s: %w", episodeID, err)
	}
	challenge, err := s.cr.GetChallengeByID(ctx, ep.VerificationChallengeID.UUID)
	if err != nil {
		return 0, fmt.Errorf("could not get challenge by id: %w", err)
	}

	numberAttempts, err := s.cr.CountAttempts(ctx, repository.CountAttemptsParams{
		UserID:    userID,
		EpisodeID: episodeID,
		CreatedAt: sql.NullTime{
			Time:  time.Now().Add(-s.activatedRealmPeriod),
			Valid: true,
		},
	})
	if err == nil {
		if numberAttempts >= int64(challenge.UserMaxAttempts) {
			return 0, nil
		}
	}

	return int64(challenge.UserMaxAttempts) - numberAttempts, nil
}
