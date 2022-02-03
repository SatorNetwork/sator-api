package result_table

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room/quiz_engine/result_table/cell"
)

type ResultTable interface {
	RegisterQuestionSendingEvent(questionNum int) error
	RegisterAnswer(userID uuid.UUID, qNum int, isCorrect bool, answeredAt time.Time) error
	GetAnswer(userID uuid.UUID, qNum int) (cell.Cell, error)
	GetPrizePoolDistribution() map[uuid.UUID]float64
	GetWinners() ([]*Winner, error)

	calcWinnersMap() map[uuid.UUID]uint32
	calcPTSMap() map[uuid.UUID]uint32
	getWinnerIDs() []uuid.UUID
}

type user struct {
	id  uuid.UUID
	pts uint32
}

type Winner struct {
	UserID string
	Prize  string
}

type Config struct {
	QuestionNum        int
	WinnersNum         int
	PrizePool          float64
	TimePerQuestionSec int
}

type resultTable struct {
	table          map[uuid.UUID][]cell.Cell
	tableMutex     *sync.Mutex
	questionSentAt []time.Time

	stakeLevels interfaces.StakeLevels

	cfg *Config
}

func New(cfg *Config, stakeLevels interfaces.StakeLevels) ResultTable {
	return &resultTable{
		table:          make(map[uuid.UUID][]cell.Cell),
		tableMutex:     &sync.Mutex{},
		questionSentAt: make([]time.Time, cfg.QuestionNum),

		stakeLevels: stakeLevels,

		cfg: cfg,
	}
}

func (rt *resultTable) RegisterQuestionSendingEvent(questionNum int) error {
	return rt.setQuestionSentAt(questionNum, time.Now())
}

func (rt *resultTable) RegisterAnswer(userID uuid.UUID, qNum int, isCorrect bool, answeredAt time.Time) error {
	rt.tableMutex.Lock()
	defer rt.tableMutex.Unlock()

	rt.createRowIfNecessary(userID)

	qSentAt, err := rt.getQuestionSentAt(qNum)
	if err != nil {
		return err
	}

	cell := cell.New(
		isCorrect,
		rt.isFirstCorrectAnswer(qNum),
		qSentAt,
		answeredAt,
		rt.cfg.TimePerQuestionSec,
	)
	return rt.setCell(userID, qNum, cell)
}

func (rt *resultTable) isFirstCorrectAnswer(qNum int) bool {
	for _, row := range rt.table {
		cell := row[qNum]
		if cell != nil && cell.IsCorrect() {
			return false
		}
	}

	return true
}

func (rt *resultTable) GetAnswer(userID uuid.UUID, qNum int) (cell.Cell, error) {
	rt.tableMutex.Lock()
	defer rt.tableMutex.Unlock()

	return rt.getCell(userID, qNum)
}

func (rt *resultTable) GetPrizePoolDistribution() map[uuid.UUID]float64 {
	rt.tableMutex.Lock()
	defer rt.tableMutex.Unlock()

	winnersMap := rt.calcWinnersMap()

	var totalPTS uint32
	for _, pts := range winnersMap {
		totalPTS += pts
	}

	distribution := make(map[uuid.UUID]float64)
	for userID, pts := range winnersMap {
		distribution[userID] = rt.cfg.PrizePool / float64(totalPTS) * float64(pts)
	}

	return distribution
}

func (rt *resultTable) GetWinners() ([]*Winner, error) {
	userIDToPrize := rt.GetPrizePoolDistribution()

	winners := make([]*Winner, 0, len(userIDToPrize))
	for userID, prize := range userIDToPrize {
		multiplier, err := rt.stakeLevels.GetMultiplier(context.Background(), userID)
		if err != nil {
			return nil, errors.Wrap(err, "could not get user's multiplier")
		}

		prize = prize / 100 * float64(multiplier)

		winners = append(winners, &Winner{
			UserID: userID.String(),
			Prize:  fmt.Sprintf("%v", prize),
		})
	}

	return winners, nil
}

func (rt *resultTable) calcWinnersMap() map[uuid.UUID]uint32 {
	ptsMap := rt.calcPTSMap()
	winnerIDs := rt.getWinnerIDs()

	winnersMap := make(map[uuid.UUID]uint32, rt.cfg.WinnersNum)
	for _, winnerID := range winnerIDs {
		winnersMap[winnerID] = ptsMap[winnerID]
	}

	return winnersMap
}

func (rt *resultTable) calcPTSMap() map[uuid.UUID]uint32 {
	ptsMap := make(map[uuid.UUID]uint32, 0)
	for userID, row := range rt.table {
		for _, cell := range row {
			ptsMap[userID] += cell.PTS()
		}
	}

	return ptsMap
}

func (rt *resultTable) getWinnerIDs() []uuid.UUID {
	users := rt.getUsersSortedByPTS()
	winnersNum := minInt(rt.cfg.WinnersNum, len(users))
	winners := users[:winnersNum]

	winnerIDs := make([]uuid.UUID, 0, len(winners))
	for _, winner := range winners {
		winnerIDs = append(winnerIDs, winner.id)
	}

	return winnerIDs
}

func (rt *resultTable) getUsersSortedByPTS() []*user {
	ptsMap := rt.calcPTSMap()

	ptsSlice := make([]*user, 0, len(ptsMap))
	for userID, pts := range ptsMap {
		ptsSlice = append(ptsSlice, &user{
			id:  userID,
			pts: pts,
		})
	}

	sort.Slice(ptsSlice, func(i, j int) bool {
		return ptsSlice[i].pts > ptsSlice[j].pts
	})

	return ptsSlice
}

func (rt *resultTable) getCell(userID uuid.UUID, qNum int) (cell.Cell, error) {
	row, ok := rt.table[userID]
	if !ok {
		return nil, NewErrRowNotFound(userID)
	}

	if qNum >= len(row) {
		return nil, NewErrCellNotFound(userID, qNum, len(row))
	}
	cell := row[qNum]

	return cell, nil
}

func (rt *resultTable) setCell(userID uuid.UUID, qNum int, cell cell.Cell) error {
	row, ok := rt.table[userID]
	if !ok {
		return NewErrRowNotFound(userID)
	}

	if qNum >= len(row) {
		return NewErrCellNotFound(userID, qNum, len(row))
	}
	row[qNum] = cell

	return nil
}

func (rt *resultTable) createRowIfNecessary(userID uuid.UUID) {
	if _, ok := rt.table[userID]; !ok {
		rt.table[userID] = make([]cell.Cell, rt.cfg.QuestionNum)
		for idx := range rt.table[userID] {
			rt.setCell(userID, idx, cell.Empty())
		}
	}
}

func (rt *resultTable) getQuestionSentAt(qNum int) (time.Time, error) {
	if qNum < 0 || qNum >= len(rt.questionSentAt) {
		return time.Time{}, NewErrIndexOutOfRange(len(rt.questionSentAt), qNum)
	}

	return rt.questionSentAt[qNum], nil
}

func (rt *resultTable) setQuestionSentAt(qNum int, qSentAt time.Time) error {
	if qNum >= len(rt.questionSentAt) {
		return NewErrIndexOutOfRange(len(rt.questionSentAt), qNum)
	}
	rt.questionSentAt[qNum] = qSentAt

	return nil
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
