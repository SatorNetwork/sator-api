package quiz

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/SatorNetwork/sator-api/svc/questions"
	"github.com/SatorNetwork/sator-api/svc/quiz/repository"

	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

// Quiz stages
const (
	WaitingForPlayers = iota << 2
	InProgress
	Finished
)

type (
	// Hub quiz hub
	Hub struct {
		Stage int

		ChallengeID          uuid.UUID
		QuizID               uuid.UUID
		PrizePool            float64
		PlayersNumberToStart int
		TotalQuestions       int
		Countdown            int
		CountdownSleep       time.Duration
		TimePerQuestion      time.Duration
		TimeBtwQuestions     time.Duration
		RewardAssetName      string

		sendMsg                 broadcast.Broadcaster
		sendQuestionResultEvent broadcast.Broadcaster
		startQuiz               broadcast.Broadcaster
		stopQuiz                broadcast.Broadcaster

		questions       map[string]questions.Question
		questionsSentAt map[string]time.Time
		answers         map[string]uuid.UUID

		players map[string]*PlayerHub
	}

	PlayerHub struct {
		UserID   uuid.UUID
		Username string

		answers map[string]QuestionResult

		sendMsg     broadcast.Broadcaster
		receiveAnsw broadcast.Broadcaster
		quit        broadcast.Broadcaster
	}
)

// NewPlayerHub setup new player hub
func NewPlayerHub(userID uuid.UUID, username string) *PlayerHub {
	return &PlayerHub{
		UserID:      userID,
		Username:    username,
		sendMsg:     broadcast.NewBroadcaster(10),
		receiveAnsw: broadcast.NewBroadcaster(10),
		quit:        broadcast.NewBroadcaster(10),
		answers:     make(map[string]QuestionResult),
	}
}

// Close player hub
func (ph *PlayerHub) Close() error {
	if err := ph.sendMsg.Close(); err != nil {
		return fmt.Errorf("could not close sending message broadcast: %w", err)
	}

	if err := ph.receiveAnsw.Close(); err != nil {
		return fmt.Errorf("could not close received answers broadcast: %w", err)
	}

	return nil
}

// NewQuizHub setup new quiz hub
func NewQuizHub(quiz repository.Quiz, qs map[string]questions.Question) *Hub {
	return &Hub{
		ChallengeID:          quiz.ChallengeID,
		QuizID:               quiz.ID,
		PrizePool:            quiz.PrizePool,
		PlayersNumberToStart: int(quiz.PlayersToStart),
		TotalQuestions:       len(qs),
		TimePerQuestion:      time.Duration(quiz.TimePerQuestion) * time.Second,
		TimeBtwQuestions:     time.Second * 2,
		Countdown:            3,
		CountdownSleep:       time.Millisecond * 700,
		RewardAssetName:      "SAO",

		sendMsg:                 broadcast.NewBroadcaster(100),
		sendQuestionResultEvent: broadcast.NewBroadcaster(100),
		startQuiz:               broadcast.NewBroadcaster(10),

		questions:       qs,
		questionsSentAt: make(map[string]time.Time),
		answers:         make(map[string]uuid.UUID),

		players: make(map[string]*PlayerHub),
	}
}

// AddPlayer adds player to a quiz hub.
func (h *Hub) AddPlayer(userID uuid.UUID, username string) error {
	if _, ok := h.players[userID.String()]; !ok {
		log.Printf("AddPlayer %s: %v", username, userID)
		if h.Stage != OpenForRegistration {
			return fmt.Errorf("quiz is closed for a new players")
		}
		if len(h.players) >= h.PlayersNumberToStart {
			if h.Stage == OpenForRegistration {
				h.SendQuizStartEvent()
			}
			return fmt.Errorf("quiz is full, try another one")
		}

		h.players[userID.String()] = NewPlayerHub(userID, username)

		if len(h.players) >= h.PlayersNumberToStart {
			h.Stage = ClosedForRegistration
			h.SendQuizStartEvent()
		}
	}

	return nil
}

// Connect adds player to a quiz hub and send other players user_connected message
func (h *Hub) Connect(userID uuid.UUID) error {
	if p, ok := h.players[userID.String()]; ok {
		log.Printf("Connect: %v", userID)
		if err := h.SendMessage(UserConnectedMessage, User{
			UserID:   userID.String(),
			Username: p.Username,
		}); err != nil {
			return fmt.Errorf("add player: could not encode message: %w", err)
		}
	} else {
		log.Printf("could not found player with id=%s", userID)
	}

	return nil
}

func (h *Hub) RemovePlayer(userID uuid.UUID) error {
	if p, ok := h.players[userID.String()]; !ok {
		if err := p.Close(); err != nil {
			return fmt.Errorf("remove player: could not close player hub: %w", err)
		}

		if err := h.SendMessage(UserDisconnectedMessage, User{
			UserID:   userID.String(),
			Username: p.Username,
		}); err != nil {
			return fmt.Errorf("remove player: could not encode message: %w", err)
		}

		delete(h.players, userID.String())

		if len(h.players) < 1 {
			h.SendQuizStopEvent()
			h.Shutdown()
		}
	}

	return nil
}

func (h *Hub) SendPlayerQuitEvent(userID uuid.UUID) {
	if _, ok := h.players[userID.String()]; ok {
		h.players[userID.String()].quit.Submit(userID.String())
	}
}

func (h *Hub) ListenPlayerQuitEvent(userID uuid.UUID, ch chan interface{}) {
	if _, ok := h.players[userID.String()]; ok {
		h.players[userID.String()].quit.Register(ch)
	}
}

func (h *Hub) UnsubscribePlayerQuitEvent(userID uuid.UUID, ch chan interface{}) {
	if _, ok := h.players[userID.String()]; ok {
		h.players[userID.String()].quit.Unregister(ch)
	}
}

func (h *Hub) ListenMessageToSend(userID uuid.UUID, ch chan interface{}) {
	h.sendMsg.Register(ch)
	if _, ok := h.players[userID.String()]; ok {
		h.players[userID.String()].sendMsg.Register(ch)
	}
}

func (h *Hub) UnsubscribeMessageToSend(userID uuid.UUID, ch chan interface{}) {
	h.sendMsg.Unregister(ch)
	if _, ok := h.players[userID.String()]; ok {
		h.players[userID.String()].sendMsg.Unregister(ch)
	}
}

// SendMessage sends message to general quiz channel
func (h *Hub) SendMessage(msgType string, msg interface{}) error {
	b, err := json.Marshal(Message{
		Type:    msgType,
		SentAt:  time.Now(),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}
	h.sendMsg.Submit(b)
	return nil
}

// SendPersonalMessage sends message to general quiz channel
func (h *Hub) SendPersonalMessage(userID uuid.UUID, msgType string, msg interface{}) error {
	b, err := json.Marshal(Message{
		Type:    msgType,
		SentAt:  time.Now(),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}

	if ph, ok := h.players[userID.String()]; ok {
		ph.sendMsg.Submit(b)
		return nil
	} else {
		log.Printf("SendPersonalMessage: could not found player hub by user_id=%s", userID.String())
	}

	return fmt.Errorf("player with id=%s not found", userID.String())
}

func (h *Hub) SendCountdownMessages() error {
	// countdown: 3.. 2.. 1.. 0!
	for i := h.Countdown; i > -1; i-- {
		if err := h.SendMessage(CountdownMessage, Countdown{
			Countdown: i,
		}); err != nil {
			return fmt.Errorf("could not sent countdown message: %w", err)
		}
		time.Sleep(h.CountdownSleep)
	}
	return nil
}

// SendQuestionResultEvent fire event that means users should receive current question result
func (h *Hub) SendQuestionResultEvent(questionID uuid.UUID) error {
	h.sendQuestionResultEvent.Submit(questionID.String())
	return nil
}

// ListenQuestionResultEvent subscribes for question result events
func (h *Hub) ListenQuestionResultEvent(ch chan interface{}) {
	h.sendQuestionResultEvent.Register(ch)
}

// UnsubscribeQuestionResultEvent unsubscribes from question result events
func (h *Hub) UnsubscribeQuestionResultEvent(ch chan interface{}) {
	h.sendQuestionResultEvent.Unregister(ch)
}

// SendQuizStartEvent fire event that means quiz should start
func (h *Hub) SendQuizStartEvent() {
	h.startQuiz.Submit(h.QuizID.String())
}

// SendQuizStopEvent fire event that means quiz should stop
func (h *Hub) SendQuizStopEvent() {
	h.stopQuiz.Submit(h.QuizID.String())
}

func (h *Hub) ListenQuizStartEvents(ch chan interface{}) {
	h.startQuiz.Register(ch)
}

func (h *Hub) UnsubscribeQuizStartEvents(ch chan interface{}) {
	h.startQuiz.Unregister(ch)
}

// SendQuestions sends quiz questions
func (h *Hub) SendQuestions() error {
	for _, q := range h.questions {
		// cast answer options
		answOpts := make([]AnswerOption, 0, len(q.AnswerOptions))
		var correctAnswID uuid.UUID
		for _, a := range q.AnswerOptions {
			answOpts = append(answOpts, AnswerOption{
				AnswerID:   a.ID.String(),
				AnswerText: a.Option,
			})
			if a.IsCorrect {
				correctAnswID = a.ID
			}
		}

		if correctAnswID == uuid.Nil {
			return fmt.Errorf("question with id=%s does not have right answer option", q.ID.String())
		}

		// send question
		if err := h.SendMessage(QuestionMessage, Question{
			QuestionID:     q.ID.String(),
			QuestionText:   q.Question,
			TimeForAnswer:  int(h.TimePerQuestion.Seconds()),
			TotalQuestions: h.TotalQuestions,
			QuestionNumber: len(h.questionsSentAt) + 1,
			AnswerOptions:  answOpts,
		}); err != nil {
			return fmt.Errorf("could not sent countdown message: %w", err)
		}
		// store when question is sent
		// it's needed to calculate answer rate and pts
		h.questionsSentAt[q.ID.String()] = time.Now()
		h.answers[q.ID.String()] = correctAnswID

		// time for sending of answer
		time.Sleep(h.TimePerQuestion)

		// fire event "send question result"
		err := h.SendQuestionResultEvent(q.ID)
		if err != nil {
			return fmt.Errorf("could not fire event sent question result: %w", err)
		}

		// time to check answers
		time.Sleep(h.TimeBtwQuestions)
	}

	return nil
}

func (h *Hub) CheckAnswer(questionID, answerID uuid.UUID) bool {
	if correctAnsw, ok := h.answers[questionID.String()]; ok && correctAnsw == answerID {
		return true
	}
	return false
}

func (h *Hub) GetCorrectAnswer(questionID uuid.UUID) uuid.UUID {
	return h.answers[questionID.String()]
}

func (h *Hub) SinceQuestionSent(questionID uuid.UUID) (float64, error) {
	if sentAt, ok := h.questionsSentAt[questionID.String()]; ok {
		return time.Since(sentAt).Seconds(), nil
	}
	return 0, fmt.Errorf("question is not sent")
}

func (h *Hub) SetAnswer(userID, questionID, answerID uuid.UUID, isCorrect bool, rate, pts int) QuestionResult {
	qr := QuestionResult{
		QuestionID:      questionID.String(),
		Result:          isCorrect,
		Rate:            rate,
		AdditionalPts:   pts,
		CorrectAnswerID: h.answers[questionID.String()].String(),
		QuestionsLeft:   h.TotalQuestions - len(h.questionsSentAt),
	}

	h.players[userID.String()].answers[questionID.String()] = qr

	return qr
}

func (h *Hub) SendQuestionResult(userID, questionID uuid.UUID) error {
	result, ok := h.players[userID.String()].answers[questionID.String()]
	if !ok {
		result = QuestionResult{
			QuestionID:      questionID.String(),
			Result:          false,
			CorrectAnswerID: h.GetCorrectAnswer(questionID).String(),
		}
	}
	if err := h.SendPersonalMessage(userID, QuestionResultMessage, result); err != nil {
		return fmt.Errorf("could not sent quiz result message: %w", err)
	}
	if !result.Result {
		time.Sleep(h.TimeBtwQuestions)
		h.SendPlayerQuitEvent(userID)
		err := h.RemovePlayer(userID)
		if err != nil {
			return fmt.Errorf("could not remove player: %w", err)
		}
	}
	return nil
}

func (h *Hub) SendWinners(winners []Winner) error {
	// sort winners in descending order by prize amount
	sort.SliceStable(winners, func(i, j int) bool {
		return winners[i].PrizeAmount > winners[j].PrizeAmount
	})

	// send quiz result
	if err := h.SendMessage(ChallengeResultMessage, ChallengeResult{
		ChallengeID: h.ChallengeID.String(),
		PrizePool:   fmt.Sprintf("%.2f %s", h.PrizePool, h.RewardAssetName),
		Winners:     winners,
	}); err != nil {
		return fmt.Errorf("could not sent quiz result message: %w", err)
	}

	return nil
}

func (h *Hub) Shutdown() {
	h.SendQuizStopEvent()
	h.startQuiz.Close()
	h.sendQuestionResultEvent.Close()
	h.sendMsg.Close()
	for _, ph := range h.players {
		ph.Close()
	}
}

func (h *Hub) notifyPlayer(uid uuid.UUID) error {
	player, ok := h.players[uid.String()]
	if !ok {
		return fmt.Errorf("could not receive player from map")
	}

	for k, v := range h.players {
		if k != uid.String() {
			b, err := json.Marshal(Message{
				Type:    UserConnectedMessage,
				SentAt:  time.Now(),
				Payload: User{
					UserID:   v.UserID.String(),
					Username: v.Username,
				},
			})
			if err != nil {
				return fmt.Errorf("could not encode message: %w", err)
			}

			player.sendMsg.Submit(b)
		}
	}

	return nil
}