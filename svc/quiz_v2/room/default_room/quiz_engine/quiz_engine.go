package quiz_engine

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	quiz_v2_challenge "github.com/SatorNetwork/sator-api/svc/quiz_v2/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table"
)

const defaultWinnersNum = 2

type resultTable interface {
	SaveUserAnswer(userID uuid.UUID, qNum int, isCorrect bool, answeredAt time.Time) (int, bool)
	GetPrizePoolDistribution() map[uuid.UUID]float64
	RegisterQuestionSendingEvent(questionNum int)
}

type QuizEngine struct {
	challengeID   string
	challengesSvc quiz_v2_challenge.ChallengesService

	challenge   *challenge.RawChallenge
	questions   []challenge.Question
	resultTable resultTable
}

func New(challengeID string, challengesSvc quiz_v2_challenge.ChallengesService) *QuizEngine {
	return &QuizEngine{
		challengeID:   challengeID,
		challengesSvc: challengesSvc,

		questions: make([]challenge.Question, 0),
	}
}

func (e *QuizEngine) Init() error {
	if err := e.loadChallenge(); err != nil {
		return err
	}
	if err := e.loadQuestions(); err != nil {
		return err
	}

	e.resultTable = result_table.New(
		len(e.questions),
		defaultWinnersNum,
		e.challenge.PrizePoolAmount,
		int(e.challenge.TimePerQuestionSec),
	)

	return nil
}

func (e *QuizEngine) loadChallenge() error {
	ctx := context.Background()
	challengeID, err := uuid.Parse(e.challengeID)
	if err != nil {
		return err
	}
	challenge, err := e.challengesSvc.GetRawChallengeByID(ctx, challengeID)
	if err != nil {
		return err
	}
	fmt.Printf("******************************************************************************************\n")
	fmt.Printf("%+v\n", challenge)
	fmt.Printf("******************************************************************************************\n\n\n")

	e.challenge = &challenge
	return nil
}

func (e *QuizEngine) loadQuestions() error {
	ctx := context.Background()
	challengeID, err := uuid.Parse(e.challengeID)
	if err != nil {
		return err
	}
	questions, err := e.challengesSvc.GetQuestionsByChallengeID(ctx, challengeID)
	if err != nil {
		return errors.Wrap(err, "can't get questions by challenge id")
	}

	for _, q := range questions {
		fmt.Printf("################################################################################\n")
		fmt.Printf("%+v\n", q)
		fmt.Printf("################################################################################\n\n\n")
	}

	e.questions = questions
	return nil
}

func (e *QuizEngine) GetQuestions() []challenge.Question {
	return e.questions
}

func (e *QuizEngine) CheckAndTrackAnswer(questionID, answerID, userID uuid.UUID, answeredAt time.Time) (bool, int, bool, error) {
	question, err := e.getQuestionByID(questionID)
	if err != nil {
		return false, 0, false, err
	}

	answer, err := e.getAnswerByID(question, answerID)
	if err != nil {
		return false, 0, false, err
	}

	segmentNum, isFastestAnswer := e.resultTable.SaveUserAnswer(userID, int(question.Order)-1, answer.IsCorrect, answeredAt)

	return answer.IsCorrect, segmentNum, isFastestAnswer, nil
}

func (e *QuizEngine) getQuestionByID(questionID uuid.UUID) (challenge.Question, error) {
	for _, q := range e.questions {
		if q.ID == questionID {
			return q, nil
		}
	}

	return challenge.Question{}, errors.Errorf("question not found")
}

func (e *QuizEngine) getAnswerByID(question challenge.Question, answerID uuid.UUID) (challenge.AnswerOption, error) {
	for _, answer := range question.AnswerOptions {
		if answer.ID == answerID {
			return answer, nil
		}
	}

	return challenge.AnswerOption{}, errors.Errorf("answer not found")
}

func (e *QuizEngine) GetPrizePoolDistribution() map[uuid.UUID]float64 {
	return e.resultTable.GetPrizePoolDistribution()
}

func (e *QuizEngine) RegisterQuestionSendingEvent(questionNum int) {
	e.resultTable.RegisterQuestionSendingEvent(questionNum)
}
