package default_room

import (
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	quiz_v2_challenge "github.com/SatorNetwork/sator-api/svc/quiz_v2/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/status_transactor"
)

const (
	defaultChanBuffSize = 10
)

type questionWrapper struct {
	question    challenge.Question
	questionNum int
}

type answerWrapper struct {
	message    *message.AnswerMessage
	receivedAt time.Time
}

type defaultRoom struct {
	challengeID    string
	players        map[string]player.Player
	playersMutex   *sync.Mutex
	newPlayersChan chan player.Player
	countdownChan  chan int
	questionChan   chan *questionWrapper
	answersChan    chan *answerWrapper

	statusIsUpdatedChan chan struct{}
	st                  *status_transactor.StatusTransactor

	quizEngine quiz_engine.QuizEngine

	done chan struct{}
}

func New(challengeID string, challenges quiz_v2_challenge.ChallengesService) (*defaultRoom, error) {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	quizEngine, err := quiz_engine.New(challengeID, challenges)
	if err != nil {
		return nil, err
	}

	return &defaultRoom{
		challengeID:    challengeID,
		players:        make(map[string]player.Player, defaultChanBuffSize),
		playersMutex:   &sync.Mutex{},
		newPlayersChan: make(chan player.Player, defaultChanBuffSize),
		countdownChan:  make(chan int, defaultChanBuffSize),
		questionChan:   make(chan *questionWrapper, defaultChanBuffSize),
		answersChan:    make(chan *answerWrapper, defaultChanBuffSize),

		statusIsUpdatedChan: statusIsUpdatedChan,
		st:                  status_transactor.New(statusIsUpdatedChan),

		quizEngine: quizEngine,

		done: make(chan struct{}),
	}, nil
}

func (r *defaultRoom) ChallengeID() string {
	return r.challengeID
}

func (r *defaultRoom) AddPlayer(p player.Player) {
	r.playersMutex.Lock()
	r.players[p.ID()] = p
	r.playersMutex.Unlock()

	r.newPlayersChan <- p
}

func (r *defaultRoom) IsFull() bool {
	challenge := r.quizEngine.GetChallenge()

	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	return len(r.players) >= int(challenge.PlayersToStart)
}

func (r *defaultRoom) Start() {
LOOP:
	for {
		select {
		case p := <-r.newPlayersChan:
			// TODO(evg): we need proper mechanism to wait until user starts to listen messages
			// accumulate messages on our side until nats connection will set up by user?
			time.Sleep(time.Second)

			r.sendPlayerIsJoinedMessage(p)
			go r.watchPlayerMessages(p)

			if r.IsFull() {
				r.st.SetStatus(status_transactor.RoomIsFullStatus)
			}

		case <-r.statusIsUpdatedChan:
			switch r.st.GetStatus() {
			case status_transactor.RoomIsFullStatus:
				// runCountdown is short-running goroutine with auto-completion, so should not be considered as a goroutine leak
				go r.runCountdown()

			case status_transactor.CountdownFinishedStatus:
				go r.runQuestions()

			case status_transactor.QuestionAreSentStatus:
				go r.sendWinnersTable()

			case status_transactor.WinnersTableAreSent:
				r.st.SetStatus(status_transactor.RoomIsFinished)

			case status_transactor.RoomIsFinished:
				go func() {
					// wait for some time to drain all messages from channels before closing the room
					time.Sleep(time.Second)
					r.Close()
				}()
			}

		case answer := <-r.answersChan:
			if err := r.processAnswerMessage(answer); err != nil {
				log.Printf("can't process answer message: %v\n", err)
				break
			}

		case secondsLeft := <-r.countdownChan:
			r.sendCountdownMessage(secondsLeft)

		case q := <-r.questionChan:
			r.sendQuestionMessage(q)
			if err := r.quizEngine.RegisterQuestionSendingEvent(q.questionNum); err != nil {
				log.Println(err)
			}

		case <-r.done:
			break LOOP
		}
	}

	log.Println("room is gracefully closed")
}

func (r *defaultRoom) Close() {
	// If room is already closed then return from function. We don't want to accidentally close `done` channel two or more times.
	if r.st.GetStatus() == status_transactor.RoomIsClosed {
		return
	}
	r.st.SetStatus(status_transactor.RoomIsClosed)

	close(r.done)
}

func (r *defaultRoom) getPlayerByID(id string) player.Player {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	return r.players[id]
}

func (r *defaultRoom) watchPlayerMessages(p player.Player) {
LOOP:
	for {
		select {
		case msg := <-p.GetMessageStream():
			r.answersChan <- &answerWrapper{
				message:    msg.MustGetAnswerMessage(),
				receivedAt: time.Now(),
			}
		case <-r.done:
			break LOOP
		}
	}
}

func (r *defaultRoom) runCountdown() {
	secondsLeft := 3
	// first message should be sent without ticker delay
	r.countdownChan <- secondsLeft
	secondsLeft--

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
LOOP:
	for ; secondsLeft >= 1; secondsLeft-- {
		select {
		case <-ticker.C:
			r.countdownChan <- secondsLeft

		case <-r.done:
			break LOOP
		}
	}
	// wait a second for last message
	time.Sleep(time.Second)

	r.st.SetStatus(status_transactor.CountdownFinishedStatus)
}

func (r *defaultRoom) runQuestions() {
	challenge := r.quizEngine.GetChallenge()
	questions := r.quizEngine.GetQuestions()

	// first question should be sent without ticker delay
	r.questionChan <- &questionWrapper{
		question:    questions[0],
		questionNum: 0,
	}

	delayBetweenQuestions := time.Duration(challenge.TimePerQuestionSec) * time.Second
	ticker := time.NewTicker(delayBetweenQuestions)
	defer ticker.Stop()
LOOP:
	for i := 1; i < len(questions); i++ {
		select {
		case <-ticker.C:
			r.questionChan <- &questionWrapper{
				question:    questions[i],
				questionNum: i,
			}

		case <-r.done:
			break LOOP
		}
	}
	// wait `delayBetweenQuestions` time to collect answers on last question
	time.Sleep(delayBetweenQuestions)

	r.st.SetStatus(status_transactor.QuestionAreSentStatus)
}

func (r *defaultRoom) sendWinnersTable() {
	challenge := r.quizEngine.GetChallenge()

	userIDToPrize := r.quizEngine.GetPrizePoolDistribution()
	usernameIDToPrize := make(map[string]float64, len(userIDToPrize))
	
	r.playersMutex.Lock()
	for userID, prize := range userIDToPrize {
		username := r.players[userID.String()].Username()
		usernameIDToPrize[username] = prize
	}

	winners := r.quizEngine.GetWinners()
	msgWinners := make([]*message.Winner, 0, len(winners))
	for _, w := range winners {
		username := r.players[w.UserID].Username()

		msgWinners = append(msgWinners, &message.Winner{
			UserID:   w.UserID,
			Username: username,
			Prize:    w.Prize,
		})
	}
	r.playersMutex.Unlock()

	payload := message.WinnersTableMessage{
		ChallengeID:           r.ChallengeID(),
		PrizePool:             challenge.PrizePool,
		ShowTransactionURL:    "TODO",
		Winners:               msgWinners,
		PrizePoolDistribution: usernameIDToPrize,
	}
	msg, err := message.NewWinnersTableMessage(&payload)
	if err != nil {
		log.Println(err)
		return
	}

	r.sendMessageToRoom(msg)

	r.st.SetStatus(status_transactor.WinnersTableAreSent)
}

func (r *defaultRoom) sendPlayerIsJoinedMessage(p player.Player) {
	payload := message.PlayerIsJoinedMessage{
		PlayerID: p.ID(),
		Username: p.Username(),
	}
	msg, err := message.NewPlayerIsJoinedMessage(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendCountdownMessage(secondsLeft int) {
	payload := message.CountdownMessage{
		SecondsLeft: secondsLeft,
	}
	msg, err := message.NewCountdownMessage(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendQuestionMessage(q *questionWrapper) {
	answerOptions := make([]message.AnswerOption, 0, len(q.question.AnswerOptions))
	for _, answer := range q.question.AnswerOptions {
		answerOptions = append(answerOptions, message.AnswerOption{
			AnswerID:   answer.ID.String(),
			AnswerText: answer.Option,
		})
	}

	payload := message.QuestionMessage{
		QuestionID:     q.question.ID.String(),
		QuestionText:   q.question.Question,
		TimeForAnswer:  q.question.TimeForAnswer,
		QuestionNumber: q.questionNum,
		AnswerOptions:  answerOptions,
	}
	msg, err := message.NewQuestionMessage(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendMessageToRoom(message *message.Message) {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	for _, p := range r.players {
		if err := p.SendMessage(message); err != nil {
			log.Printf("can't send message to player with %v uid: %v\n", p.ID(), err)
		}
	}
}

func (r *defaultRoom) processAnswerMessage(answer *answerWrapper) error {
	userID, err := uuid.Parse(answer.message.UserID)
	if err != nil {
		return errors.Wrapf(err, "can't parse user's UID(%v)", answer.message.UserID)
	}
	questionID, err := uuid.Parse(answer.message.QuestionID)
	if err != nil {
		return errors.Wrapf(err, "can't parse question's UID(%v)", answer.message.QuestionID)
	}
	answerID, err := uuid.Parse(answer.message.AnswerID)
	if err != nil {
		return errors.Wrapf(err, "can't parse answer's UID(%v)", answer.message.AnswerID)
	}

	ok, err := r.quizEngine.CheckAndRegisterAnswer(questionID, answerID, userID, answer.receivedAt)
	if err != nil {
		return errors.Wrapf(err, "can't check answer, question UID(%v), answer's UID(%v)", questionID, answerID)
	}
	cell, err := r.quizEngine.GetAnswer(userID, questionID)
	if err != nil {
		return err
	}

	answerID, err = r.quizEngine.GetCorrectAnswerID(questionID)
	if err != nil {
		return err
	}
	questionNum, err := r.quizEngine.GetQuestionNumByID(questionID)
	if err != nil {
		return err
	}
	questionsLeft := r.quizEngine.GetNumberOfQuestions() - questionNum - 1

	payload := message.AnswerReplyMessage{
		QuestionID:      questionID.String(),
		Success:         ok,
		Rate:            cell.Rate(),
		CorrectAnswerID: answerID.String(),
		QuestionsLeft:   questionsLeft,
		AdditionalPTS:   cell.AdditionalPTS(),
		SegmentNum:      cell.FindSegmentNum(),
		IsFastestAnswer: cell.IsFirstCorrectAnswer(),
	}
	msg, err := message.NewAnswerReplyMessage(&payload)
	if err != nil {
		return err
	}

	p := r.getPlayerByID(answer.message.UserID)
	if err := p.SendMessage(msg); err != nil {
		return errors.Wrapf(err, "can't send message to player with %v uid", answer.message.UserID)
	}

	return nil
}