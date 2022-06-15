package quiz_engine

import (
	"time"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/common"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/question_container"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table/cell"
)

type QuizEngine interface {
	GetCurrentPrizePool() float64
	GetChallenge() *challenge.RawChallenge
	GetNumberOfQuestions() int
	GetQuestions() []challenge.Question
	GetQuestionNumByID(questionID uuid.UUID) (int, error)
	GetCorrectAnswerID(questionID uuid.UUID) (uuid.UUID, error)
	CheckAndRegisterAnswer(questionID, answerID, userID uuid.UUID, answeredAt time.Time) (bool, error)
	GetAnswer(userID, questionID uuid.UUID) (cell.Cell, error)
	RegisterQuestionSendingEvent(questionNum int) error
	GetPrizePoolDistribution() (map[uuid.UUID]result_table.UserReward, error)
	GetPlayers() ([]*result_table.Player, error)
	GetWinnersAndLosers() ([]*result_table.Winner, []*result_table.Loser, error)
	GetWinners() ([]*result_table.Winner, error)
	DistributedPrizePool() float64
}

type quizEngine struct {
	questionContainer question_container.QuestionContainer
	resultTable       result_table.ResultTable

	currentPrizePool float64
}

func New(
	challengeID string,
	challengesSvc interfaces.ChallengesService,
	qr interfaces.QuizV2Repository,
	stakeLevels interfaces.StakeLevels,
	shuffleQuestions bool,
) (*quizEngine, error) {
	qc, err := question_container.New(challengeID, challengesSvc, shuffleQuestions)
	if err != nil {
		return nil, err
	}
	challenge := qc.GetChallenge()

	currentPrizePool, err := getCurrentPrizePool(qc, qr)
	if err != nil {
		return nil, err
	}

	cfg := result_table.Config{
		QuestionNum:        qc.GetNumberOfQuestions(),
		WinnersNum:         int(challenge.MaxWinners),
		PrizePool:          currentPrizePool,
		TimePerQuestionSec: int(challenge.TimePerQuestionSec),
		MinCorrectAnswers:  challenge.MinCorrectAnswers,
	}
	rt := result_table.New(&cfg, stakeLevels)

	return &quizEngine{
		questionContainer: qc,
		resultTable:       rt,
		currentPrizePool:  currentPrizePool,
	}, nil
}

func getCurrentPrizePool(qc question_container.QuestionContainer, qr interfaces.QuizV2Repository) (float64, error) {
	challenge := qc.GetChallenge()
	return common.GetCurrentPrizePool(
		qr,
		challenge.ID,
		challenge.PrizePoolAmount,
		challenge.MinimumReward,
		challenge.PercentForQuiz,
	)
}

func (e *quizEngine) GetCurrentPrizePool() float64 {
	return e.currentPrizePool
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

func (e *quizEngine) GetPrizePoolDistribution() (map[uuid.UUID]result_table.UserReward, error) {
	return e.resultTable.GetPrizePoolDistribution()
}

func (e *quizEngine) GetPlayers() ([]*result_table.Player, error) {
	return e.resultTable.GetPlayers()
}

func (e *quizEngine) GetWinnersAndLosers() ([]*result_table.Winner, []*result_table.Loser, error) {
	return e.resultTable.GetWinnersAndLosers()
}

func (e *quizEngine) GetWinners() ([]*result_table.Winner, error) {
	return e.resultTable.GetWinners()
}

func (e *quizEngine) DistributedPrizePool() float64 {
	return e.resultTable.DistributedPrizePool()
}
