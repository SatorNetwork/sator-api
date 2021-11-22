package default_room

import (
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/status_transactor"
)

const (
	defaultChanBuffSize   = 10
	defaultPlayersNum     = 2
	delayBetweenQuestions = 100 * time.Millisecond
)

type Question struct {
	Text string
}

type Answer struct {
	AnswerFlag bool
	UserID     string
}

var defaultQuestions = []*Question{
	{
		Text: "question1",
	},
	{
		Text: "question2",
	},
	{
		Text: "question3",
	},
}

type defaultRoom struct {
	challengeID    string
	players        map[string]player.Player
	newPlayersChan chan player.Player
	countdownChan  chan int
	questionChan   chan *Question
	answersChan    chan *Answer

	statusIsUpdatedChan chan struct{}
	st                  *status_transactor.StatusTransactor

	done chan struct{}
}

func New(challengeID string) *defaultRoom {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	return &defaultRoom{
		challengeID:    challengeID,
		players:        make(map[string]player.Player, defaultChanBuffSize),
		newPlayersChan: make(chan player.Player, defaultChanBuffSize),
		countdownChan:  make(chan int, defaultChanBuffSize),
		questionChan:   make(chan *Question, defaultChanBuffSize),
		answersChan:    make(chan *Answer, defaultChanBuffSize),

		statusIsUpdatedChan: statusIsUpdatedChan,
		st:                  status_transactor.New(statusIsUpdatedChan),

		done: make(chan struct{}),
	}
}

func (r *defaultRoom) ChallengeID() string {
	return r.challengeID
}

func (r *defaultRoom) AddPlayer(p player.Player) {
	r.newPlayersChan <- p
}

func (r *defaultRoom) AnswerToQuestion(a *Answer) {
	r.answersChan <- a
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
				r.st.SetStatus(status_transactor.RoomIsFinished)

			case status_transactor.RoomIsFinished:
				go func() {
					// wait for some time to drain all messages from channels before closing the room
					time.Sleep(time.Second)
					r.Close()
				}()
			}

		case answer := <-r.answersChan:
			payload := message.AnswerReplyMessage{
				Success: answer.AnswerFlag == true,
			}
			msg := message.NewAnswerReplyMessage(&payload)

			if err := r.players[answer.UserID].SendMessage(msg); err != nil {
				log.Printf("can't send message to player with %v uid: %v\n", answer.UserID, err)
			}

		case secondsLeft := <-r.countdownChan:
			r.sendCountdownMessage(secondsLeft)

		case q := <-r.questionChan:
			r.sendQuestionMessage(q)

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
			r.answersChan <- &Answer{
				AnswerFlag: msg.AnswerMessage.AnswerFlag,
				UserID:     msg.AnswerMessage.UserID,
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
	// first question should be sent without ticker delay
	r.questionChan <- defaultQuestions[0]

	ticker := time.NewTicker(delayBetweenQuestions)
	defer ticker.Stop()
LOOP:
	for i := 1; i < len(defaultQuestions); i++ {
		q := defaultQuestions[i]
		select {
		case <-ticker.C:
			r.questionChan <- q

		case <-r.done:
			break LOOP
		}
	}

	r.st.SetStatus(status_transactor.QuestionAreSentStatus)
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

func (r *defaultRoom) sendQuestionMessage(q *Question) {
	payload := message.QuestionMessage{
		Text: q.Text,
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
