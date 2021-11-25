package result_table

import (
	"time"
)

const segmentsNum = 4

type answerStat struct {
	isCorrect            bool
	qSentAt              time.Time
	answeredAt           time.Time
	timePerQuestionSec   int
	isFirstCorrectAnswer bool
}

func (st *answerStat) pts() uint32 {
	if !st.isCorrect {
		return 0
	}

	return 1 + st.extraPts()
}

func (st *answerStat) extraPts() uint32 {
	pts := st.extraPtsForSegment()
	if st.isFirstCorrectAnswer {
		pts += 2
	}

	return pts
}

func (st *answerStat) extraPtsForSegment() uint32 {
	segmentNum := st.findSegmentNum()

	return map[int]uint32{
		1: 3,
		2: 2,
		3: 1,
		4: 0,
	}[segmentNum]
}

func (st *answerStat) findSegmentNum() int {
	answerTime := st.answerTime()
	segmentDuration := float64(st.timePerQuestionSec) / segmentsNum

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

	//fmt.Printf("******************************************************************************************\n")
	//fmt.Printf("ANSWER_TIME      : %+v\n", answerTime)
	//fmt.Printf("SEGMENT_DURATION : %+v\n", segmentDuration)
	//fmt.Printf("SEGMENT_NUM      : %+v\n", segmentNum)
	//fmt.Printf("******************************************************************************************\n\n\n")

	return segmentNum
}

func (st *answerStat) answerTime() time.Duration {
	//fmt.Printf("******************************************************************************************\n")
	//fmt.Printf("ANSWERED_AT : %+v\n", st.answeredAt)
	//fmt.Printf("Q_SENT_AT   : %+v\n", st.qSentAt)
	//fmt.Printf("******************************************************************************************\n\n\n")

	return st.answeredAt.Sub(st.qSentAt)
}
