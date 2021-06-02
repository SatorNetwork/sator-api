package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/questions"
	"github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/dmitrymomot/go-signature"
	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		mutex          mutex
		repo           quizRepository
		questions      questionService
		rewards        rewardsService
		challenges     challengesService
		tokenGenFunc   tokenGenFunc
		tokenParseFunc tokenParseFunc
		tokenTTL       int64
		baseQuizURL    string

		countdown       int
		timeForAnswer   time.Duration
		timeBtwQuestion time.Duration
		rewardAssetName string

		startQuiz broadcast.Broadcaster // receives quiz id to start
		stopQuiz  broadcast.Broadcaster // receives quiz id to stop
		quizzes   map[string]quizHub
		players   map[string]playerHub
	}

	playerHub struct {
		send    broadcast.Broadcaster
		receive broadcast.Broadcaster
	}

	quizHub struct {
		send            broadcast.Broadcaster
		questResult     broadcast.Broadcaster // fire event to send question result to each user
		questionsSentAt map[string]time.Time  // time when question was sent
		questionsNumber int
	}

	quizRepository interface {
		AddNewQuiz(ctx context.Context, arg repository.AddNewQuizParams) (repository.Quiz, error)
		GetQuizByChallengeID(ctx context.Context, challengeID uuid.UUID) (repository.Quiz, error)
		GetQuizByID(ctx context.Context, id uuid.UUID) (repository.Quiz, error)
		GetQuizWinnners(ctx context.Context, arg repository.GetQuizWinnnersParams) ([]repository.GetQuizWinnnersRow, error)
		UpdateQuizStatus(ctx context.Context, arg repository.UpdateQuizStatusParams) error
		GetAnswer(ctx context.Context, arg repository.GetAnswerParams) (repository.QuizAnswer, error)
		AddNewPlayer(ctx context.Context, arg repository.AddNewPlayerParams) error
		StoreAnswer(ctx context.Context, arg repository.StoreAnswerParams) error
	}

	challengesService interface {
		GetChallengeByID(ctx context.Context, challengeID uuid.UUID) (challenge.Challenge, error)
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
		UserID   string
		Username string
		QuizID   string
	}

	ServiceOption func(s *Service)

	mutex interface {
		Lock(key string, ttl time.Duration) error
		Unlock(key string) error
	}

	PlayerAnswer struct {
		QuizID     string `json:"quiz_id"`
		UserID     string `json:"user_id"`
		QuestionID string `json:"question_id"`
		AnswerID   string `json:"answer_id"`
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
		mutex:           m,
		repo:            repo,
		questions:       questions,
		rewards:         rewards,
		challenges:      challenges,
		tokenGenFunc:    signature.NewTemporary,
		tokenParseFunc:  signature.Parse,
		tokenTTL:        900,
		baseQuizURL:     baseQuizURL,
		countdown:       3,
		timeForAnswer:   time.Second * 15,
		timeBtwQuestion: time.Second * 5,
		rewardAssetName: "SAO",

		startQuiz: broadcast.NewBroadcaster(100),
		stopQuiz:  broadcast.NewBroadcaster(100), // receives quiz id to stop
		quizzes:   make(map[string]quizHub),
		players:   make(map[string]playerHub),
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

func (s *Service) Play(ctx context.Context, quizID, uid uuid.UUID, username string) error {
	qh := s.GetQuizHub(quizID.String())
	ph := s.GetPlayerHub(uid.String())

	s.sendQuizMsg(quizID.String(), UserConnectedMessage, User{
		UserID:   uid.String(),
		Username: username,
	})

	questResult := make(chan interface{}, 100)
	qh.questResult.Register(questResult)
	playerSend := make(chan interface{}, 100)
	ph.send.Register(playerSend)
	playerAnswers := make(chan interface{}, 100)
	ph.receive.Register(playerAnswers)

	defer func() {
		qh.questResult.Unregister(questResult)
		close(questResult)
		ph.send.Unregister(playerSend)
		close(playerSend)
		ph.receive.Unregister(playerAnswers)
		close(playerAnswers)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("quiz service stopped")
			return nil

		case msg := <-questResult:
			if questID, ok := msg.(string); ok {
				result := false
				answ, err := s.repo.GetAnswer(ctx, repository.GetAnswerParams{
					QuizID:     quizID,
					UserID:     uid,
					QuestionID: uuid.MustParse(questID),
				})
				if err == nil && answ.IsCorrect {
					result = true
				}

				if err := s.sendPersonalMsg(uid.String(), QuestionResultMessage, QuestionResult{
					QuestionID:    questID,
					Result:        result,
					Rate:          int(answ.Rate),
					AdditionalPts: int(answ.Pts),
					QuestionsLeft: qh.questionsNumber - len(qh.questionsSentAt),
				}); err != nil {
					return fmt.Errorf("could not sent challenge result message for question with id=%s: %w", questID, err)
				}

				if !result {
					log.Printf("user %s lose in quiz %s", uid.String(), quizID.String())
					return nil
				}
			}

		case answ := <-playerAnswers:
			if a, ok := answ.(PlayerAnswer); ok {
				answID := uuid.MustParse(a.AnswerID)
				correct, err := s.questions.CheckAnswer(ctx, answID)
				if err != nil {
					correct = false
				}
				rate := int32(s.calcRate(a.QuizID, a.QuestionID, int64(s.timeForAnswer.Seconds())))
				pts := 0
				if rate > 2 {
					pts = 2
				}
				if err := s.repo.StoreAnswer(ctx, repository.StoreAnswerParams{
					QuizID:     uuid.MustParse(a.QuizID),
					UserID:     uuid.MustParse(a.UserID),
					QuestionID: uuid.MustParse(a.QuestionID),
					AnswerID:   answID,
					IsCorrect:  correct,
					Rate:       rate,
					Pts:        int32(pts),
				}); err != nil {
					log.Printf("could not store player`s (%s) answer (%s) for question (%s): %v",
						uid.String(), a.AnswerID, a.QuestionID, err)
				}
			}
		}
	}
}

func (s *Service) calcRate(quizID, questID string, timeForAnswer int64) int {
	qh := s.GetQuizHub(quizID)
	if t, ok := qh.questionsSentAt[questID]; ok {
		answTime := int64(time.Since(t).Seconds())
		partDur := timeForAnswer / 4
		if answTime <= partDur {
			return 3
		} else if answTime <= (partDur * 2) {
			return 2
		} else if answTime <= (partDur * 3) {
			return 1
		}
	}
	return 0
}

func (s *Service) SubmitAnswer(quizID, uid string, answ Answer) {
	ph := s.GetPlayerHub(uid)
	ph.receive.Submit(PlayerAnswer{
		QuizID:     quizID,
		UserID:     uid,
		QuestionID: answ.QuestionID,
		AnswerID:   answ.AnswerID,
	})
}

func (s *Service) GetQuizHub(qid string) quizHub {
	qb, ok := s.quizzes[qid]
	if !ok {
		s.quizzes[qid] = quizHub{
			send:            broadcast.NewBroadcaster(100),
			questResult:     broadcast.NewBroadcaster(100),
			questionsSentAt: make(map[string]time.Time),
		}
	}
	return qb
}

func (s *Service) GetPlayerHub(uid string) playerHub {
	br, ok := s.players[uid]
	if !ok {
		s.players[uid] = playerHub{
			send:    broadcast.NewBroadcaster(10),
			receive: broadcast.NewBroadcaster(10),
		}
	}
	return br
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
					} else {
						log.Printf("quiz with id %s has been finished", q.ID.String())
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
	for i := s.countdown; i > -1; i-- {
		if err := s.sendQuizMsg(quiz.ID.String(), CountdownMessage, Countdown{
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
		if err := s.sendQuizMsg(q.ID.String(), QuestionMessage, Question{
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
	if err := s.sendQuizMsg(quiz.ID.String(), ChallengeResultMessage, ChallengeResult{
		ChallengeID: quiz.ID.String(),
		PrizePool:   fmt.Sprintf("%v %s", quiz.PrizePool, s.rewardAssetName),
		Winners:     winners,
	}); err != nil {
		return fmt.Errorf("could not sent challenge result message for quiz with id=%s: %w", quiz.ID.String(), err)
	}
	// wait for question timer
	time.Sleep(time.Duration(quiz.TimePerQuestion) * time.Second)

	s.stopQuiz.Submit(quiz)

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

	for _, w := range result {
		totalPts += int(w.Pts)
	}

	for _, w := range result {
		prize := calcPrize(quiz.PrizePool, totalWinnersNumber, int(questionsNumber), totalPts, int(w.Pts))
		if err := s.rewards.AddReward(ctx, w.UserID, prize, quiz.ID); err != nil {
			log.Printf("could not store reward: user_id=%s, quiz_id=%s, amount=%v error: %v",
				w.UserID.String(), quiz.ID.String(), prize, err)
		}
		winners = append(winners, Winner{
			UserID:   w.UserID.String(),
			Username: w.Username,
			Prize:    fmt.Sprintf("%v %s", prize, s.rewardAssetName),
		})
	}

	return winners, nil
}

func calcPrize(prizePool float64, totalWinners, totalQuestions, totalPts, pts int) float64 {
	if totalWinners == 1 {
		return prizePool
	}

	totalPoints := (totalWinners * totalQuestions) + totalPts
	winnerPoints := totalQuestions + pts

	return (prizePool / float64(totalPoints)) * float64(winnerPoints)
}

func (s *Service) sendQuizMsg(qid string, msgType string, msg interface{}) error {
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

func (s *Service) sendPersonalMsg(uid string, msgType string, msg interface{}) error {
	br, ok := s.players[uid]
	if !ok {
		return fmt.Errorf("could not found player with id: %s", uid)
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
