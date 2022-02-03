package quiz_engine

import (
	"time"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/question_container"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table/cell"
)

const defaultWinnersNum = 2

type QuizEngine interface {
	GetChallenge() *challenge.RawChallenge
	GetNumberOfQuestions() int
	GetQuestions() []challenge.Question
	GetQuestionNumByID(questionID uuid.UUID) (int, error)
	GetCorrectAnswerID(questionID uuid.UUID) (uuid.UUID, error)
	CheckAndRegisterAnswer(questionID, answerID, userID uuid.UUID, answeredAt time.Time) (bool, error)
	GetAnswer(userID, questionID uuid.UUID) (cell.Cell, error)
	RegisterQuestionSendingEvent(questionNum int) error
	GetPrizePoolDistribution() map[uuid.UUID]float64
	GetWinners() ([]*result_table.Winner, error)
}

type quizEngine struct {
	questionContainer question_container.QuestionContainer
	resultTable       result_table.ResultTable
}

func New(challengeID string, challengesSvc interfaces.ChallengesService, stakeLevels interfaces.StakeLevels) (*quizEngine, error) {
	qc, err := question_container.New(challengeID, challengesSvc)
	if err != nil {
		return nil, err
	}

	cfg := result_table.Config{
		QuestionNum:        qc.GetNumberOfQuestions(),
		WinnersNum:         defaultWinnersNum,
		PrizePool:          qc.GetChallenge().PrizePoolAmount,
		TimePerQuestionSec: int(qc.GetChallenge().TimePerQuestionSec),
	}
	rt := result_table.New(&cfg, stakeLevels)

	return &quizEngine{
		questionContainer: qc,
		resultTable:       rt,
	}, nil
}

func (e *quizEngine) GetChallenge() *challenge.RawChallenge {
	return e.questionContainer.GetChallenge()
}

func (e *quizEngine) GetNumberOfQuestions() int {
	return e.questionContainer.GetNumberOfQuestions()
}

func (e *quizEngine) GetQuestions() []challenge.Question {
	return e.questionContainer.GetQuestions()
}

func (e *quizEngine) GetQuestionNumByID(questionID uuid.UUID) (int, error) {
	return e.questionContainer.GetQuestionNumByID(questionID)
}

func (e *quizEngine) GetCorrectAnswerID(questionID uuid.UUID) (uuid.UUID, error) {
	return e.questionContainer.GetCorrectAnswerID(questionID)
}

func (e *quizEngine) CheckAndRegisterAnswer(questionID, answerID, userID uuid.UUID, answeredAt time.Time) (bool, error) {
	isCorrect, err := e.questionContainer.CheckAnswer(questionID, answerID)
	if err != nil {
		return false, err
	}
	qNum, err := e.questionContainer.GetQuestionNumByID(questionID)
	if err != nil {
		return false, err
	}

	if err := e.resultTable.RegisterAnswer(userID, qNum, isCorrect, answeredAt); err != nil {
		return false, err
	}

	return isCorrect, nil
}

func (e *quizEngine) GetAnswer(userID, questionID uuid.UUID) (cell.Cell, error) {
	qNum, err := e.questionContainer.GetQuestionNumByID(questionID)
	if err != nil {
		return nil, err
	}

	return e.resultTable.GetAnswer(userID, qNum)
}

func (e *quizEngine) RegisterQuestionSendingEvent(questionNum int) error {
	return e.resultTable.RegisterQuestionSendingEvent(questionNum)
}

func (e *quizEngine) GetPrizePoolDistribution() map[uuid.UUID]float64 {
	return e.resultTable.GetPrizePoolDistribution()
}

func (e *quizEngine) GetWinners() ([]*result_table.Winner, error) {
	return e.resultTable.GetWinners()
}
