package default_room

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/restriction_manager"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/status_transactor"
)

const (
	defaultChanBuffSize = 10

	RelationTypeQuizzes = "quizzes"
)

type questionWrapper struct {
	question    challenge.Question
	questionNum int
}

type answerWrapper struct {
	message    *message.AnswerMessage
	userID     string
	receivedAt time.Time
}

type defaultRoom struct {
	challengeID           string
	players               map[string]player.Player
	playersMutex          *sync.Mutex
	newPlayersChan        chan player.Player
	countdownChan         chan int
	questionChan          chan *questionWrapper
	answersChan           chan *answerWrapper
	messagesForNewPlayers []*message.Message

	statusIsUpdatedChan chan struct{}
	st                  *status_transactor.StatusTransactor

	quizEngine         quiz_engine.QuizEngine
	rewards            interfaces.RewardsService
	restrictionManager restriction_manager.RestrictionManager

	done chan struct{}
}

func New(
	challengeID string,
	challenges interfaces.ChallengesService,
	stakeLevels interfaces.StakeLevels,
	rewards interfaces.RewardsService,
	restrictionManager restriction_manager.RestrictionManager,
	shuffleQuestions bool,
) (*defaultRoom, error) {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	quizEngine, err := quiz_engine.New(challengeID, challenges, stakeLevels, shuffleQuestions)
	if err != nil {
		return nil, err
	}

	return &defaultRoom{
		challengeID:           challengeID,
		players:               make(map[string]player.Player, defaultChanBuffSize),
		playersMutex:          &sync.Mutex{},
		newPlayersChan:        make(chan player.Player, defaultChanBuffSize),
		countdownChan:         make(chan int, defaultChanBuffSize),
		questionChan:          make(chan *questionWrapper, defaultChanBuffSize),
		answersChan:           make(chan *answerWrapper, defaultChanBuffSize),
		messagesForNewPlayers: make([]*message.Message, 0),

		statusIsUpdatedChan: statusIsUpdatedChan,
		st:                  status_transactor.New(statusIsUpdatedChan),

		quizEngine:         quizEngine,
		rewards:            rewards,
		restrictionManager: restrictionManager,

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

// NOTE: should be run as a goroutine
func (r *defaultRoom) Start() {
LOOP:
	for {
		select {
		case p := <-r.newPlayersChan:
			{
				tmpl := `
				New player is joined
				User's ID: %v
				Username:  %v
				`
				log.Printf(tmpl, p.ID(), p.Username())
			}

			r.sendMessagesForNewPlayers(p)
			r.sendPlayerIsJoinedMessage(p)

			go r.watchPlayerMessages(p)

			if r.IsFull() {
				log.Println("Room is full")

				r.st.SetStatus(status_transactor.RoomIsFullStatus)
			}

		case <-r.statusIsUpdatedChan:
			switch r.st.GetStatus() {
			case status_transactor.RoomIsFullStatus:
				// runCountdown is short-running goroutine with auto-completion, so should not be considered as a goroutine leak
				go r.runCountdown()

			case status_transactor.CountdownFinishedStatus:
				go r.registerAttempts()
				go r.runQuestions()

			case status_transactor.QuestionAreSentStatus:
				go r.sendWinnersTable()

			case status_transactor.WinnersTableAreSent:
				go r.sendRewards()

			case status_transactor.RewardsAreSent:
				r.st.SetStatus(status_transactor.RoomIsFinished)

			case status_transactor.RoomIsFinished:
				go func() {
					// wait for some time to drain all messages from channels before closing the room
					time.Sleep(time.Second)
					r.Close()

					r.closePlayers()
				}()
			}

		case answer := <-r.answersChan:
			if err := r.processAnswerMessage(answer); err != nil {
				log.Printf("can't process answer message: %v\n", err)
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

func (r *defaultRoom) GetRoomDetails() *room.RoomDetails {
	challenge := r.quizEngine.GetChallenge()

	return &room.RoomDetails{
		PlayersToStart: challenge.PlayersToStart,
	}
}

func (r *defaultRoom) Close() {
	// If room is already closed then return from function. We don't want to accidentally close `done` channel two or more times.
	if r.st.GetStatus() == status_transactor.RoomIsClosed {
		return
	}
	r.st.SetStatus(status_transactor.RoomIsClosed)

	close(r.done)
}

func (r *defaultRoom) closePlayers() {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	for _, p := range r.players {
		if err := p.Close(); err != nil {
			log.Printf("can't close player: %v\n", err)
		}
	}
}

func (r *defaultRoom) getPlayerByID(id string) player.Player {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	return r.players[id]
}

func (r *defaultRoom) getPlayerByIDNoLock(id string) player.Player {
	return r.players[id]
}

func (r *defaultRoom) removePlayerByID(id string) {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	delete(r.players, id)
}

func (r *defaultRoom) watchPlayerMessages(p player.Player) {
LOOP:
	for {
		select {
		case msg := <-p.GetMessageStream():
			r.answersChan <- &answerWrapper{
				message:    msg.MustGetAnswerMessage(),
				userID:     p.ID(),
				receivedAt: time.Now(),
			}

		case <-p.ConnectChan():
			// TODO(evg): handle userIsConnected event

		case <-p.DisconnectChan():
			if r.st.GetStatus() == status_transactor.GatheringPlayersStatus {
				r.removePlayerByID(p.ID())
				if err := p.Close(); err != nil {
					log.Printf("can't close player: %v\n", err)
				}
				r.sendPlayerIsDisconnectedMessage(p)
				break LOOP
			}

		case <-r.done:
			break LOOP
		}
	}
}

func (r *defaultRoom) runCountdown() {
	time.Sleep(time.Second * 2)

	secondsLeft := 3
	// first message should be sent without ticker delay
	r.countdownChan <- secondsLeft
	secondsLeft--

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
LOOP:
	for ; secondsLeft >= 0; secondsLeft-- {
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

func (r *defaultRoom) registerAttempts() {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	challengeID, err := uuid.Parse(r.ChallengeID())
	if err != nil {
		log.Printf("can't parse challenge ID: %v\n", err)
		return
	}

	for _, p := range r.players {
		playerID, err := uuid.Parse(p.ID())
		if err != nil {
			log.Printf("can't parse player ID: %v\n", err)
			continue
		}

		err = r.restrictionManager.RegisterAttempt(context.Background(), challengeID, playerID)
		if err != nil {
			log.Printf("can't register challenge attempt: %v\n", err)
			continue
		}
	}
}

func (r *defaultRoom) runQuestions() {
	challenge := r.quizEngine.GetChallenge()
	questions := r.quizEngine.GetQuestions()

	// if there are no questions - return from the function with correct status and move on
	// it should NOT happen - just precautionary measure to avoid panics (in case challenge isn't properly setup)
	if len(questions) == 0 {
		r.st.SetStatus(status_transactor.QuestionAreSentStatus)
		return
	}

	// first question should be sent without ticker delay
	r.questionChan <- &questionWrapper{
		question:    questions[0],
		questionNum: 0,
	}

	afterAnswerReplyDelay := 2 * time.Second
	delayBetweenQuestions := time.Duration(challenge.TimePerQuestionSec) * time.Second
	ticker := time.NewTicker(delayBetweenQuestions + afterAnswerReplyDelay)
	defer ticker.Stop()
LOOP:
	for i := 1; i < len(questions); i++ {
		select {
		case <-ticker.C:
			if err := r.sendAnswerReplyMessages(questions[i-1].ID); err != nil {
				log.Printf("can't send answer reply messages: %v\n", err)
			}
			time.Sleep(afterAnswerReplyDelay)

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
	qNum := len(questions)
	if err := r.sendAnswerReplyMessages(questions[qNum-1].ID); err != nil {
		log.Printf("can't send answer reply messages: %v\n", err)
	}
	time.Sleep(afterAnswerReplyDelay)

	r.st.SetStatus(status_transactor.QuestionAreSentStatus)
}

func (r *defaultRoom) sendWinnersTable() {
	challenge := r.quizEngine.GetChallenge()

	userIDToPrize, err := r.quizEngine.GetPrizePoolDistribution()
	if err != nil {
		log.Printf("can't get prize pool distribution: %v\n", err)
		return
	}
	usernameIDToPrize := make(map[string]float64, len(userIDToPrize))

	r.playersMutex.Lock()
	for userID, prize := range userIDToPrize {
		username := r.players[userID.String()].Username()
		usernameIDToPrize[username] = prize
	}

	winners, losers, err := r.quizEngine.GetWinnersAndLosers()
	if err != nil {
		log.Printf("can't get winners: %v\n", err)
		return
	}
	if int32(len(winners)) > challenge.MaxWinners {
		log.Printf(
			"max amount of winners exceeded, max winners: %v, winners: %v\n",
			challenge.MaxWinners,
			len(winners),
		)
		return
	}

	msgWinners := make([]*message.Winner, 0, len(winners))
	for _, w := range winners {
		p := r.players[w.UserID]

		msgWinners = append(msgWinners, &message.Winner{
			UserID:   w.UserID,
			Username: p.Username(),
			Prize:    w.Prize,
			Avatar:   p.Avatar(),
		})
	}

	msgLosers := make([]*message.Loser, 0, len(losers))
	for _, loser := range losers {
		p := r.players[loser.UserID]

		msgLosers = append(msgLosers, &message.Loser{
			UserID:   loser.UserID,
			Username: p.Username(),
			PTS:      loser.PTS,
			Avatar:   p.Avatar(),
		})
	}

	r.playersMutex.Unlock()

	payload := message.WinnersTableMessage{
		ChallengeID:           r.ChallengeID(),
		PrizePool:             challenge.PrizePool,
		ShowTransactionURL:    "TODO",
		Winners:               msgWinners,
		Losers:                msgLosers,
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

func (r *defaultRoom) sendRewards() {
	userIDToPrize, err := r.quizEngine.GetPrizePoolDistribution()
	if err != nil {
		log.Printf("can't get prize pool distribution: %v\n", err)
		return
	}
	challengeID, err := uuid.Parse(r.ChallengeID())
	if err != nil {
		log.Printf("can't parse challenge ID: %v\n", err)
		return
	}

	for userID, prize := range userIDToPrize {
		err := r.rewards.AddDepositTransaction(
			context.Background(),
			userID,
			challengeID,
			RelationTypeQuizzes,
			prize,
		)
		if err != nil {
			log.Printf("can't add deposit transaction: %v\n", err)
			continue
		}
	}

	for userID, prize := range userIDToPrize {
		err := r.restrictionManager.RegisterEarnedReward(context.Background(), challengeID, userID, prize)
		if err != nil {
			log.Printf("can't register earned reward: %v\n", err)
			return
		}
	}

	r.st.SetStatus(status_transactor.RewardsAreSent)
}

func (r *defaultRoom) sendMessagesForNewPlayers(p player.Player) {
	for _, msg := range r.messagesForNewPlayers {
		if err := p.SendMessage(msg); err != nil {
			log.Printf("can't send message: %v\n", err)
		}
	}
}

func (r *defaultRoom) sendPlayerIsJoinedMessage(p player.Player) {
	payload := message.PlayerIsJoinedMessage{
		PlayerID: p.ID(),
		Username: p.Username(),
		Avatar:   p.Avatar(),
	}
	msg, err := message.NewPlayerIsJoinedMessage(&payload)
	if err != nil {
		log.Println(err)
		return
	}
	r.messagesForNewPlayers = append(r.messagesForNewPlayers, msg)
	r.sendMessageToRoom(msg)
}

func (r *defaultRoom) sendPlayerIsDisconnectedMessage(p player.Player) {
	payload := message.PlayerIsDisconnectedMessage{
		PlayerID: p.ID(),
		Username: p.Username(),
	}
	msg, err := message.NewPlayerIsDisconnectedMessage(&payload)
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
	challenge := r.quizEngine.GetChallenge()

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
		TimeForAnswer:  int(challenge.TimePerQuestionSec),
		QuestionNumber: q.questionNum,
		TotalQuestions: r.quizEngine.GetNumberOfQuestions(),
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
	userID, err := uuid.Parse(answer.userID)
	if err != nil {
		return errors.Wrapf(err, "can't parse user's UID(%v)", answer.userID)
	}
	questionID, err := uuid.Parse(answer.message.QuestionID)
	if err != nil {
		return errors.Wrapf(err, "can't parse question's UID(%v)", answer.message.QuestionID)
	}
	answerID, err := uuid.Parse(answer.message.AnswerID)
	if err != nil {
		return errors.Wrapf(err, "can't parse answer's UID(%v)", answer.message.AnswerID)
	}

	_, err = r.quizEngine.CheckAndRegisterAnswer(questionID, answerID, userID, answer.receivedAt)
	if err != nil {
		return errors.Wrapf(err, "can't check answer, question UID(%v), answer's UID(%v)", questionID, answerID)
	}

	return nil
}

func (r *defaultRoom) sendAnswerReplyMessages(questionID uuid.UUID) error {
	r.playersMutex.Lock()
	defer r.playersMutex.Unlock()

	for _, p := range r.players {
		playerID, err := uuid.Parse(p.ID())
		if err != nil {
			log.Printf("can't parse player uid: %v\n", err)
			continue
		}

		err = r.sendAnswerReplyMessage(playerID, questionID)
		if err != nil {
			log.Printf("can't sendAnswerReplyMessage: %v\n", err)
			continue
		}
	}

	return nil
}

func (r *defaultRoom) sendAnswerReplyMessage(userID, questionID uuid.UUID) error {
	cell, err := r.quizEngine.GetAnswer(userID, questionID)
	if err != nil {
		return err
	}

	answerID, err := r.quizEngine.GetCorrectAnswerID(questionID)
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
		Success:         cell.IsCorrect(),
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

	p := r.getPlayerByIDNoLock(userID.String())
	if err := p.SendMessage(msg); err != nil {
		return errors.Wrapf(err, "can't send message to player with %v uid", userID.String())
	}

	return nil
}
