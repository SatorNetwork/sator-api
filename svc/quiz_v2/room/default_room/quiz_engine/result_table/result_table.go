package result_table

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type user struct {
	id  uuid.UUID
	pts uint32
}

type ResultTable struct {
	table              map[uuid.UUID][]*answerStat
	questionSentAt     []time.Time
	questionNum        int
	winnersNum         int
	prizePool          float64
	timePerQuestionSec int
}

func New(questionNum, winnersNum int, prizePool float64, timePerQuestionSec int) *ResultTable {
	return &ResultTable{
		table:              make(map[uuid.UUID][]*answerStat),
		questionSentAt:     make([]time.Time, questionNum),
		questionNum:        questionNum,
		winnersNum:         winnersNum,
		prizePool:          prizePool,
		timePerQuestionSec: timePerQuestionSec,
	}
}

func (rt *ResultTable) SaveUserAnswer(userID uuid.UUID, qNum int, isCorrect bool, answeredAt time.Time) (int, bool) {
	if _, ok := rt.table[userID]; !ok {
		rt.table[userID] = make([]*answerStat, rt.questionNum)
	}
	row := rt.table[userID]
	qSentAt := rt.questionSentAt[qNum]

	row[qNum] = &answerStat{
		isCorrect:            isCorrect,
		qSentAt:              qSentAt,
		answeredAt:           answeredAt,
		timePerQuestionSec:   rt.timePerQuestionSec,
		isFirstCorrectAnswer: rt.isFirstCorrectAnswer(qNum),
	}

	cell := row[qNum]
	return cell.findSegmentNum(), cell.isFirstCorrectAnswer
}

func (rt *ResultTable) isFirstCorrectAnswer(qNum int) bool {
	for _, row := range rt.table {
		cell := row[qNum]
		if cell != nil && cell.isCorrect {
			return false
		}
	}

	return true
}

func (rt *ResultTable) GetPrizePoolDistribution() map[uuid.UUID]float64 {
	winnersMap := rt.calcWinnersMap()
	fmt.Printf("******************************************************************************************\n")
	fmt.Printf("%+v\n", winnersMap)
	fmt.Printf("******************************************************************************************\n\n\n")

	var totalPTS uint32
	for _, pts := range winnersMap {
		totalPTS += pts
	}

	distribution := make(map[uuid.UUID]float64)
	for userID, pts := range winnersMap {
		distribution[userID] = rt.prizePool / float64(totalPTS) * float64(pts)
	}

	return distribution
}

func (rt *ResultTable) calcWinnersMap() map[uuid.UUID]uint32 {
	ptsMap := rt.calcPTSMap()
	winnerIDs := rt.getWinnerIDs()

	winnersMap := make(map[uuid.UUID]uint32, rt.winnersNum)
	for _, winnerID := range winnerIDs {
		winnersMap[winnerID] = ptsMap[winnerID]
	}

	return winnersMap
}

func (rt *ResultTable) calcPTSMap() map[uuid.UUID]uint32 {
	ptsMap := make(map[uuid.UUID]uint32, 0)
	for userID, row := range rt.table {
		for _, answerStat := range row {
			ptsMap[userID] += answerStat.pts()
		}
	}

	return ptsMap
}

func (rt *ResultTable) getWinnerIDs() []uuid.UUID {
	users := rt.getUsersSortedByPTS()
	winnersNum := minInt(rt.winnersNum, len(users))
	winners := users[:winnersNum]

	winnerIDs := make([]uuid.UUID, 0, len(winners))
	for _, winner := range winners {
		winnerIDs = append(winnerIDs, winner.id)
	}

	return winnerIDs
}

func (rt *ResultTable) getUsersSortedByPTS() []*user {
	ptsMap := rt.calcPTSMap()

	ptsSlice := make([]*user, 0, len(ptsMap))
	for userID, pts := range ptsMap {
		ptsSlice = append(ptsSlice, &user{
			id:  userID,
			pts: pts,
		})
	}

	sort.Slice(ptsSlice, func(i, j int) bool {
		return ptsSlice[i].pts < ptsSlice[j].pts
	})

	return ptsSlice
}

func (rt *ResultTable) RegisterQuestionSendingEvent(questionNum int) {
	rt.questionSentAt[questionNum] = time.Now()
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
