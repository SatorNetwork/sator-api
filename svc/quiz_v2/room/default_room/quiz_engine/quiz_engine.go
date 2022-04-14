package quiz_engine

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/question_container"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table/cell"
)

type QuizEngine interface {
	GetChallenge() *challenge.RawChallenge
	GetNumberOfQuestions() int
	GetQuestions() []challenge.Question
	GetQuestionNumByID(questionID uuid.UUID) (int, error)
	GetCorrectAnswerID(questionID uuid.UUID) (uuid.UUID, error)
	CheckAndRegisterAnswer(questionID, answerID, userID uuid.UUID, answeredAt time.Time) (bool, error)
	GetAnswer(userID, questionID uuid.UUID) (cell.Cell, error)
	RegisterQuestionSendingEvent(questionNum int) error
	GetPrizePoolDistribution() (map[uuid.UUID]result_table.UserReward, error)
	GetWinnersAndLosers() ([]*result_table.Winner, []*result_table.Loser, error)
	GetWinners() ([]*result_table.Winner, error)
	DistributedPrizePool() float64
}

type quizEngine struct {
	questionContainer question_container.QuestionContainer
	resultTable       result_table.ResultTable
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
	}, nil
}

func getCurrentPrizePool(qc question_container.QuestionContainer, qr interfaces.QuizV2Repository) (float64, error) {
	challenge := qc.GetChallenge()

	ctxb := context.Background()
	distributedRewards, err := qr.GetDistributedRewardsByChallengeID(ctxb, qc.GetChallenge().ID)
	if err != nil && !strings.Contains(err.Error(), "converting NULL to float64 is unsupported") {
		return 0, errors.Wrap(err, "can't get distributed rewards by challenge id")
	}
	if err != nil && strings.Contains(err.Error(), "converting NULL to float64 is unsupported") {
		distributedRewards = 0
	}
	leftInPool := challenge.PrizePoolAmount - distributedRewards
	if leftInPool <= 0 {
		return 0, errors.Wrap(err, "no money left in pool")
	}
	if leftInPool <= challenge.MinimumReward {
		return leftInPool, nil
	}

	currentPrizePool := leftInPool / 100 * challenge.PercentForQuiz
	if currentPrizePool < challenge.MinimumReward {
		currentPrizePool = challenge.MinimumReward
	}

	return currentPrizePool, nil

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

func (e *quizEngine) GetWinnersAndLosers() ([]*result_table.Winner, []*result_table.Loser, error) {
	return e.resultTable.GetWinnersAndLosers()
}

func (e *quizEngine) GetWinners() ([]*result_table.Winner, error) {
	return e.resultTable.GetWinners()
}

func (e *quizEngine) DistributedPrizePool() float64 {
	return e.resultTable.DistributedPrizePool()
}
