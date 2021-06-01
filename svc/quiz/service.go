package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/questions"
	"github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/dmitrymomot/go-signature"
	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		repo           quizRepository
		questions      questionService
		rewards        rewardsService
		tokenGenFunc   tokenGenFunc
		tokenParseFunc tokenParseFunc
		tokenTTL       int64
		baseQuizURL    string

		countdown       int
		timeForAnswer   time.Duration
		timeBtwQuestion time.Duration

		startQuiz broadcast.Broadcaster // receives quiz id to start
		stopQuiz  broadcast.Broadcaster // receives quiz id to stop
		quizzes   map[string]struct {
			send        broadcast.Broadcaster
			questResult broadcast.Broadcaster // fire event to send question result to each user
		}
	}

	quizRepository interface {
		AddNewQuiz(ctx context.Context, arg repository.AddNewQuizParams) (repository.Quiz, error)
		GetQuizByChallengeID(ctx context.Context, challengeID uuid.UUID) (repository.Quiz, error)
		GetQuizByID(ctx context.Context, id uuid.UUID) (repository.Quiz, error)
		GetQuizWinnners(ctx context.Context, arg repository.GetQuizWinnnersParams) ([]repository.GetQuizWinnnersRow, error)
		UpdateQuizStatus(ctx context.Context, arg repository.UpdateQuizStatusParams) error
	}

	questionService interface {
		GetQuestionsByChallengeID(ctx context.Context, challengeID uuid.UUID) ([]questions.Question, error)
		CheckAnswer(ctx context.Context, answerID uuid.UUID) (bool, error)
	}

	rewardsService interface {
		AddReward(ctx context.Context, userID uuid.UUID, amount float64, quizID uuid.UUID) error
	}

	tokenGenFunc   func(data interface{}, ttl int64) (string, error)
	tokenParseFunc func(token string) (interface{}, error)

	TokenPayload struct {
		UserID          string
		Username        string
		ChallengeRoomID string
	}

	ServiceOption func(s *Service)
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo quizRepository, baseQuizURL string, opt ...ServiceOption) *Service {
	s := &Service{
		repo:            repo,
		tokenGenFunc:    signature.NewTemporary,
		tokenParseFunc:  signature.Parse,
		tokenTTL:        900,
		baseQuizURL:     baseQuizURL,
		countdown:       3,
		timeForAnswer:   time.Second * 15,
		timeBtwQuestion: time.Second * 5,
	}

	for _, fn := range opt {
		fn(s)
	}

	return s
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
	startCh := make(chan interface{}, 100)
	s.startQuiz = broadcast.NewBroadcaster(1000)
	s.startQuiz.Register(startCh)

	stopCh := make(chan interface{}, 100)
	s.stopQuiz = broadcast.NewBroadcaster(1000)
	s.stopQuiz.Register(stopCh)

	defer func() {
		s.startQuiz.Unregister(startCh)
		s.startQuiz.Close()
		close(startCh)

		s.stopQuiz.Unregister(stopCh)
		s.stopQuiz.Close()
		close(stopCh)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("quiz service stopped")
			return nil

		case q := <-startCh:
			if qz, ok := q.(repository.Quiz); ok {
				go func(c context.Context, q repository.Quiz) {
					if err := s.runQuiz(c, q); err != nil {
						log.Printf("quiz with id %s is finished with error: %v", q.ID.String(), err)
					}
				}(ctx, qz)
			}

		case q := <-stopCh:
			if qz, ok := q.(repository.Quiz); ok {
				if qch, ok := s.quizzes[qz.ID.String()]; ok {
					qch.send.Close()
					qch.questResult.Close()
				}
			}
		}
	}
}

func (s *Service) runQuiz(ctx context.Context, quiz repository.Quiz) error {
	// countdown: 3.. 2.. 1.. 0!
	for i := s.countdown; i > 0; i-- {
		if err := s.sendMsg(quiz.ID.String(), CountdownMessage, Countdown{
			Countdown: i,
		}); err != nil {
			return fmt.Errorf("could not sent countdown message: %w", err)
		}
		time.Sleep(time.Second)
	}

	questions, err := s.questions.GetQuestionsByChallengeID(ctx, quiz.ChallengeID)
	if err != nil {
		return fmt.Errorf("could not get questions for quiz with id=%s: %w", quiz.ID.String(), err)
	}

	totalQuestions := len(questions)
	for n, q := range questions {
		// send question
		if err := s.sendMsg(q.ID.String(), QuestionMessage, Question{
			QuestionID:     q.ID.String(),
			QuestionText:   q.Question,
			TimeForAnswer:  int(quiz.TimePerQuestion),
			TotalQuestions: totalQuestions,
			QuestionNumber: n + 1,
			AnswerOptions:  castAnswerOptions(q.AnswerOptions),
		}); err != nil {
			return fmt.Errorf("could not sent countdown message: %w", err)
		}
		// wait for question timer
		time.Sleep(time.Duration(quiz.TimePerQuestion) * time.Second)

		// send event, which means that question time is over and the result can be sent
		if err := s.fireEventToSendQuestionResult(quiz.ID.String(), q.ID.String()); err != nil {
			return fmt.Errorf("could not send result for question with id=%s: %w", q.ID.String(), err)
		}
		// pause to display question result for each users
		time.Sleep(s.timeBtwQuestion)
	}

	winners, err := s.getQuizWinners(ctx, quiz, int32(totalQuestions))
	if err != nil {
		return fmt.Errorf("could not show winners for quiz with id=%s: %w", quiz.ID.String(), err)
	}

	// send question
	if err := s.sendMsg(quiz.ID.String(), ChallengeResultMessage, ChallengeResult{
		ChallengeID: quiz.ID.String(),
		PrizePool:   fmt.Sprintf("%v SAO", quiz.PrizePool),
		Winners:     winners,
	}); err != nil {
		return fmt.Errorf("could not sent challenge result message for quiz with id=%s: %w", quiz.ID.String(), err)
	}
	// wait for question timer
	time.Sleep(time.Duration(quiz.TimePerQuestion) * time.Second)

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
	for _, w := range result {
		winners = append(winners, Winner{
			UserID:   w.UserID.String(),
			Username: w.Username,
			Prize:    fmt.Sprintf("%v SOA", s.calcPrize(quiz.PrizePool, w.Pts, totalWinnersNumber)),
		})
	}

	return winners, nil
}

func (s *Service) calcPrize(prizePool float64, pts int32, totalWinners int) float64 {
	return 0
}

func (s *Service) sendMsg(qid string, msgType string, msg interface{}) error {
	br, ok := s.quizzes[qid]
	if !ok {
		return fmt.Errorf("could not found quiz with id: %s", qid)
	}

	b, err := json.Marshal(Message{
		Type:    msgType,
		SentAt:  time.Now(),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}

	br.send.Submit(b)

	return nil
}

func (s *Service) fireEventToSendQuestionResult(quizID, questionID string) error {
	br, ok := s.quizzes[quizID]
	if !ok {
		return fmt.Errorf("could not found quiz with id: %s", quizID)
	}

	br.questResult.Submit(questionID)

	return nil
}

func castAnswerOptions(source []questions.AnswerOption) []AnswerOption {
	result := make([]AnswerOption, 0, len(source))
	for _, a := range source {
		result = append(result, AnswerOption{
			AnswerID:   a.ID.String(),
			AnswerText: a.Option,
		})
	}
	return result
}
