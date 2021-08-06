package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/questions"
	"github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/dmitrymomot/go-signature"
	"github.com/google/uuid"
)

// Predefined quiz statuses
const (
	OpenForRegistration = iota << 2
	ClosedForRegistration

	// RelationTypeQuizzes indicates that relation type is "quizzes".
	RelationTypeQuizzes = "quizzes"
)

type (
	// Service struct
	Service struct {
		mutex           mutex
		repo            quizRepository
		questions       questionService
		rewards         rewardsService
		challenges      challengesService
		tokenGenFunc    tokenGenFunc
		tokenParseFunc  tokenParseFunc
		tokenTTL        int64
		baseQuizURL     string
		rewardAssetName string

		hub map[string]*Hub

		startQuizEvent chan interface{}
		stopQuizEvent  chan interface{}

		numberOfQuestions int
	}

	quizRepository interface {
		AddNewQuiz(ctx context.Context, arg repository.AddNewQuizParams) (repository.Quiz, error)
		GetQuizByChallengeID(ctx context.Context, challengeID uuid.UUID) (repository.Quiz, error)
		GetQuizByID(ctx context.Context, id uuid.UUID) (repository.Quiz, error)
		GetQuizWinnners(ctx context.Context, arg repository.GetQuizWinnnersParams) ([]repository.GetQuizWinnnersRow, error)
		UpdateQuizStatus(ctx context.Context, arg repository.UpdateQuizStatusParams) error
		GetAnswer(ctx context.Context, arg repository.GetAnswerParams) (repository.QuizAnswer, error)
		AddNewPlayer(ctx context.Context, arg repository.AddNewPlayerParams) error
		CountPlayersInQuiz(ctx context.Context, quizID uuid.UUID) (int64, error)
		StoreAnswer(ctx context.Context, arg repository.StoreAnswerParams) (repository.QuizAnswer, error)
	}

	challengesService interface {
		GetChallengeByID(ctx context.Context, challengeID uuid.UUID) (challenge.Challenge, error)
	}

	questionService interface {
		GetQuestionsByChallengeID(ctx context.Context, challengeID uuid.UUID) ([]questions.Question, error)
		CheckAnswer(ctx context.Context, answerID uuid.UUID) (bool, error)
	}

	rewardsService interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
	}

	tokenGenFunc   func(data interface{}, ttl int64) (string, error)
	tokenParseFunc func(token string) (interface{}, error)

	TokenPayload struct {
		UserID   string
		Username string
		QuizID   string
	}

	ServiceOption func(s *Service)

	mutex interface {
		Lock(key string, ttl time.Duration) error
		Unlock(key string) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(
	m mutex,
	repo quizRepository,
	questions questionService,
	rewards rewardsService,
	challenges challengesService,
	baseQuizURL string,
	opt ...ServiceOption,
) *Service {

	s := &Service{
		mutex:             m,
		repo:              repo,
		questions:         questions,
		rewards:           rewards,
		challenges:        challenges,
		tokenGenFunc:      signature.NewTemporary,
		tokenParseFunc:    signature.Parse,
		tokenTTL:          900,
		baseQuizURL:       baseQuizURL,
		rewardAssetName:   "SAO",
		numberOfQuestions: 5,

		hub: make(map[string]*Hub),
	}

	for _, fn := range opt {
		fn(s)
	}

	return s
}

// GetQuizLink returns link with token to connect to quiz
func (s *Service) GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (interface{}, error) {
	key := fmt.Sprintf("quiz_link_%s", challengeID)
	s.mutex.Lock(key, time.Second*3)
	defer s.mutex.Unlock(key)

	quiz, err := s.repo.GetQuizByChallengeID(ctx, challengeID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get quiz: %w", err)
		}

		challenge, err := s.challenges.GetChallengeByID(ctx, challengeID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge with id=%s: %w", challengeID.String(), err)
		}
		// Create new quiz if not exist
		quiz, err = s.repo.AddNewQuiz(ctx, repository.AddNewQuizParams{
			ChallengeID:     challenge.ID,
			PrizePool:       challenge.PrizePoolAmount,
			PlayersToStart:  int32(challenge.Players),
			TimePerQuestion: int64(challenge.TimePerQuestionSec),
		})
		if err != nil {
			return nil, fmt.Errorf("could not start new quiz: %w", err)
		}
	}

	playersNumber, err := s.repo.CountPlayersInQuiz(ctx, quiz.ID)
	if err != nil && !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could ont count players in current quiz: %w", err)
	}
	if playersNumber >= int64(quiz.PlayersToStart) {
		// Close quiz for registration
		if err := s.repo.UpdateQuizStatus(ctx, repository.UpdateQuizStatusParams{
			ID:     quiz.ID,
			Status: ClosedForRegistration,
		}); err != nil {
			log.Printf("could not update status of quiz with id=%s: %v", quiz.ID.String(), err)
		}

		challenge, err := s.challenges.GetChallengeByID(ctx, challengeID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge with id=%s: %w", challengeID.String(), err)
		}
		// Create new quiz
		quiz, err = s.repo.AddNewQuiz(ctx, repository.AddNewQuizParams{
			ChallengeID:     challenge.ID,
			PrizePool:       challenge.PrizePoolAmount,
			PlayersToStart:  int32(challenge.Players),
			TimePerQuestion: int64(challenge.TimePerQuestionSec),
		})
		if err != nil {
			return nil, fmt.Errorf("could not start new quiz: %w", err)
		}
	}

	if err := s.repo.AddNewPlayer(ctx, repository.AddNewPlayerParams{
		QuizID:   quiz.ID,
		UserID:   uid,
		Username: username,
	}); err != nil {
		return nil, fmt.Errorf("could not add a new player into quiz: %w", err)
	}

	token, err := s.tokenGenFunc(TokenPayload{
		UserID:   uid.String(),
		Username: username,
		QuizID:   quiz.ID.String(),
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

func (s *Service) GetQuizHub(qid string) *Hub {
	return s.hub[qid]
}

func (s *Service) SetupNewQuizHub(ctx context.Context, qid uuid.UUID) (*Hub, error) {
	if h, ok := s.hub[qid.String()]; ok {
		return h, nil
	}

	quiz, err := s.repo.GetQuizByID(ctx, qid)
	if err != nil {
		return nil, fmt.Errorf("could not get quiz with id=%s: %w", qid.String(), err)
	}

	qlist, err := s.questions.GetQuestionsByChallengeID(ctx, quiz.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("could not get questions list for quiz with id=%s: %w", qid.String(), err)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(qlist), func(i, j int) { qlist[i], qlist[j] = qlist[j], qlist[i] })

	ql := qlist
	if len(qlist) > s.numberOfQuestions {
		ql = qlist[:s.numberOfQuestions]
	}

	qlmap := make(map[string]questions.Question)
	for _, item := range ql {
		qlmap[item.ID.String()] = item
	}

	s.hub[qid.String()] = NewQuizHub(quiz, qlmap)
	s.hub[qid.String()].ListenQuizStartEvents(s.startQuizEvent)

	return s.hub[qid.String()], nil
}

func (s *Service) StoreAnswer(ctx context.Context, userID, quizID, questionID, answerID uuid.UUID) error {
	h := s.GetQuizHub(quizID.String())
	if h == nil {
		return fmt.Errorf("could not found quiz with id=%s", quizID.String())
	}

	isCorrect, err := s.questions.CheckAnswer(ctx, answerID)
	if err != nil {
		return fmt.Errorf("could not check answer: %w", err)
	}

	since, err := h.SinceQuestionSent(questionID)
	if err != nil {
		return err
	}
	rate := calcRate(h.TimePerQuestion.Seconds(), since)
	a, err := s.repo.StoreAnswer(ctx, repository.StoreAnswerParams{
		QuizID:     quizID,
		UserID:     userID,
		QuestionID: questionID,
		AnswerID:   answerID,
		IsCorrect:  isCorrect,
		Rate:       int32(rate),
	})
	if err != nil {
		return fmt.Errorf("could not store answer: %w", err)
	}

	h.SetAnswer(userID, questionID, answerID, isCorrect, int(a.Rate), int(a.Pts))

	return nil
}

func (s *Service) Play(ctx context.Context, quizID, uid uuid.UUID, username string) error {
	h := s.GetQuizHub(quizID.String())
	questionResultEvents := make(chan interface{}, 10)
	stopQuiz := make(chan interface{}, 10)
	h.ListenQuestionResultEvent(questionResultEvents)
	h.ListenQuizStartEvents(stopQuiz)

	defer func() {
		h.UnsubscribeQuestionResultEvent(questionResultEvents)
		h.UnsubscribeQuizStartEvents(stopQuiz)
		close(questionResultEvents)
		close(stopQuiz)
	}()

stop:
	for {
		select {
		case e := <-stopQuiz:
			if _, ok := e.(bool); ok && e.(bool) {
				break stop
			}
		case questionID := <-questionResultEvents:
			if err := h.SendQuestionResult(uid, uuid.MustParse(fmt.Sprintf("%v", questionID))); err != nil {
				return fmt.Errorf("could not send question result: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) Serve(ctx context.Context) error {
	s.startQuizEvent = make(chan interface{}, 100)
	s.stopQuizEvent = make(chan interface{}, 100)

	defer func() {
		close(s.startQuizEvent)
		close(s.stopQuizEvent)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("quiz service stopped")
			return nil

		case q := <-s.startQuizEvent:
			if qid, ok := q.(string); ok {
				go func(c context.Context, id string) {
					if err := s.runQuiz(c, id); err != nil {
						log.Printf("quiz with id %s is finished with error: %v", id, err)
					} else {
						log.Printf("quiz with id %s has been finished", id)
					}
				}(ctx, qid)
			}

		case q := <-s.stopQuizEvent:
			if qid, ok := q.(string); ok {
				s.stopQuizRunner(qid)
			}
		}
	}
}

func (s *Service) stopQuizRunner(quizID string) {
	if q, ok := s.hub[quizID]; ok {
		q.Shutdown()
	}
	delete(s.hub, quizID)
}

func (s *Service) runQuiz(ctx context.Context, quizID string) error {
	h, ok := s.hub[quizID]
	if !ok {
		return fmt.Errorf("could not found quiz hub with id=%s", quizID)
	}

	if err := s.repo.UpdateQuizStatus(ctx, repository.UpdateQuizStatusParams{
		ID:     h.QuizID,
		Status: ClosedForRegistration,
	}); err != nil {
		return fmt.Errorf("could not update quiz status: %w", err)
	}

	if err := h.SendCountdownMessages(); err != nil {
		return fmt.Errorf("run quiz: counntdown: %w", err)
	}

	if err := h.SendQuestions(); err != nil {
		return fmt.Errorf("run quiz: questions: %w", err)
	}

	quiz, err := s.repo.GetQuizByID(ctx, h.QuizID)
	if err != nil {
		return fmt.Errorf("could not get quiz with id=%s: %w", quizID, err)
	}
	winners, err := s.getQuizWinners(ctx, quiz, int32(h.TotalQuestions))
	if err != nil {
		return fmt.Errorf("could not show winners for quiz with id=%s: %w", quiz.ID.String(), err)
	}

	if err := h.SendWinners(winners); err != nil {
		return fmt.Errorf("run quiz: questions: %w", err)
	}

	if err := s.repo.UpdateQuizStatus(ctx, repository.UpdateQuizStatusParams{
		ID:     quiz.ID,
		Status: Finished,
	}); err != nil {
		return fmt.Errorf("could not update quiz status: %w", err)
	}

	return nil
}

func (s *Service) getQuizWinners(ctx context.Context, quiz repository.Quiz, questionsNumber int32) ([]Winner, error) {
	result, err := s.repo.GetQuizWinnners(ctx, repository.GetQuizWinnnersParams{
		QuizID:         quiz.ID,
		CorrectAnswers: questionsNumber,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return []Winner{}, nil
		}
		return nil, fmt.Errorf("could not get winners for quiz with id=%s: %w", quiz.ID.String(), err)
	}

	totalWinnersNumber := len(result)
	winners := make([]Winner, 0, totalWinnersNumber)
	totalPts := 0
	totalRate := 0

	for _, w := range result {
		totalPts += int(w.Pts)
		totalRate += int(w.Rate)
	}

	for _, w := range result {
		prize := calcPrize(quiz.PrizePool, totalWinnersNumber, int(questionsNumber), totalPts, totalRate, int(w.Pts), int(w.Rate))
		if err := s.rewards.AddDepositTransaction(ctx, w.UserID, quiz.ID, RelationTypeQuizzes, prize); err != nil {
			log.Printf("could not store reward: user_id=%s, quiz_id=%s, amount=%v error: %v",
				w.UserID.String(), quiz.ID.String(), prize, err)
		}
		winners = append(winners, Winner{
			UserID:      w.UserID.String(),
			Username:    w.Username,
			Prize:       fmt.Sprintf("%.2f %s", prize, s.rewardAssetName),
			PrizeAmount: prize,
		})
	}

	return winners, nil
}

func calcPrize(prizePool float64, totalWinners, totalQuestions, totalPts, totalRate, pts, rate int) float64 {
	if totalWinners == 1 {
		return prizePool
	}

	totalPoints := (totalWinners * totalQuestions) + totalPts + totalRate
	winnerPoints := totalQuestions + pts + rate

	return (prizePool / float64(totalPoints)) * float64(winnerPoints)
}
