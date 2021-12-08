package cell

import (
	"time"
)

const segmentsNum = 4

type Cell interface {
	PTS() uint32
	IsCorrect() bool
	IsFirstCorrectAnswer() bool
	FindSegmentNum() int
}

func New(
	isCorrect bool,
	isFirstCorrectAnswer bool,
	qSentAt time.Time,
	answeredAt time.Time,
	timePerQuestionSec int,
) Cell {
	return &cell{
		isCorrect:            isCorrect,
		isFirstCorrectAnswer: isFirstCorrectAnswer,
		qSentAt:              qSentAt,
		answeredAt:           answeredAt,
		timePerQuestionSec:   timePerQuestionSec,
	}
}

func Empty() Cell {
	return new(cell)
}

type cell struct {
	isCorrect            bool
	isFirstCorrectAnswer bool
	qSentAt              time.Time
	answeredAt           time.Time
	timePerQuestionSec   int
}

func (c *cell) PTS() uint32 {
	if !c.isCorrect {
		return 0
	}

	return 1 + c.extraPts()
}

func (c *cell) IsCorrect() bool {
	return c.isCorrect
}

func (c *cell) IsFirstCorrectAnswer() bool {
	return c.isFirstCorrectAnswer
}

func (c *cell) FindSegmentNum() int {
	answerTime := c.answerTime()
	segmentDuration := float64(c.timePerQuestionSec) / segmentsNum

	var segmentNum int
	switch {
	case answerTime.Seconds() < segmentDuration:
		segmentNum = 1
	case answerTime.Seconds() < 2*segmentDuration:
		segmentNum = 2
	case answerTime.Seconds() < 3*segmentDuration:
		segmentNum = 3
	default:
		segmentNum = 4
	}

	return segmentNum
}

func (c *cell) extraPts() uint32 {
	pts := c.extraPtsForSegment()
	if c.isFirstCorrectAnswer {
		pts += 2
	}

	return pts
}

func (c *cell) extraPtsForSegment() uint32 {
	segmentNum := c.FindSegmentNum()

	return map[int]uint32{
		1: 3,
		2: 2,
		3: 1,
		4: 0,
	}[segmentNum]
}

func (c *cell) answerTime() time.Duration {
	return c.answeredAt.Sub(c.qSentAt)
}
