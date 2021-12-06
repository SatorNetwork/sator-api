package default_room

import (
	"log"
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
	defaultChanBuffSize   = 10
	defaultPlayersNum     = 2
	delayBetweenQuestions = 100 * time.Millisecond
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
	r.newPlayersChan <- p
}

func (r *defaultRoom) IsFull() bool {
	return len(r.players) >= defaultPlayersNum
}

func (r *defaultRoom) Start() {
LOOP:
	for {
		select {
		case p := <-r.newPlayersChan:
			// TODO(evg): we need proper mechanism to wait until user starts to listen messages
			// accumulate messages on our side until nats connection will set up by user?
			time.Sleep(time.Second)

			r.players[p.ID()] = p
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
			userID, err := uuid.Parse(answer.message.UserID)
			if err != nil {
				log.Printf("can't parse user's UID(%v): %v\n", answer.message.UserID, err)
				break
			}
			questionID, err := uuid.Parse(answer.message.QuestionID)
			if err != nil {
				log.Printf("can't parse question's UID(%v): %v\n", answer.message.QuestionID, err)
				break
			}
			answerID, err := uuid.Parse(answer.message.AnswerID)
			if err != nil {
				log.Printf("can't parse answer's UID(%v): %v\n", answer.message.AnswerID, err)
				break
			}
			ok, err := r.quizEngine.CheckAndRegisterAnswer(questionID, answerID, userID, answer.receivedAt)
			if err != nil {
				log.Printf("can't check answer, question UID(%v), answer's UID(%v): %v\n", questionID, answerID, err)
				break
			}
			cell, err := r.quizEngine.GetAnswer(userID, questionID)
			if err != nil {
				log.Println(err)
				break
			}

			payload := message.AnswerReplyMessage{
				Success:         ok,
				SegmentNum:      cell.FindSegmentNum(),
				IsFastestAnswer: cell.IsFirstCorrectAnswer(),
			}
			msg := message.NewAnswerReplyMessage(&payload)

			if err := r.players[answer.message.UserID].SendMessage(msg); err != nil {
				log.Printf("can't send message to player with %v uid: %v\n", answer.message.UserID, err)
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
	close(r.done)
}

func (r *defaultRoom) watchPlayerMessages(p player.Player) {
LOOP:
	for {
		select {
		case msg := <-p.GetMessageStream():
			r.answersChan <- &answerWrapper{
				message:    msg.AnswerMessage,
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
	questions := r.quizEngine.GetQuestions()

	// first question should be sent without ticker delay
	r.questionChan <- &questionWrapper{
		question:    questions[0],
		questionNum: 0,
	}

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
	userIDToPrize := r.quizEngine.GetPrizePoolDistribution()
	usernameIDToPrize := make(map[string]float64, len(userIDToPrize))
	for userID, prize := range userIDToPrize {
		username := r.players[userID.String()].Username()
		usernameIDToPrize[username] = prize
	}

	payload := message.WinnersTableMessage{
		PrizePoolDistribution: usernameIDToPrize,
	}
	msg := message.NewWinnersTableMessage(&payload)
	r.sendMessageToRoom(msg)

	r.st.SetStatus(status_transactor.WinnersTableAreSent)
}

func (r *defaultRoom) sendPlayerIsJoinedMessage(p player.Player) {
	payload := message.PlayerIsJoinedMessage{
		PlayerID: p.ID(),
		Username: p.Username(),
	}
	msg := message.NewPlayerIsJoinedMessage(&payload)
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendCountdownMessage(secondsLeft int) {
	payload := message.CountdownMessage{
		SecondsLeft: secondsLeft,
	}
	msg := message.NewCountdownMessage(&payload)
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
	msg := message.NewQuestionMessage(&payload)
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendMessageToRoom(message *message.Message) {
	for _, p := range r.players {
		if err := p.SendMessage(message); err != nil {
			log.Printf("can't send message to player with %v uid: %v\n", p.ID(), err)
		}
	}
}