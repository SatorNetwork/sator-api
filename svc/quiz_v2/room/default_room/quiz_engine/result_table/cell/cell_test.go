package cell

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCell(t *testing.T) {
	{
		var (
			isCorrect            = true
			isFirstCorrectAnswer = true
			qSentAt              = time.Now()
			answeredAt           = time.Now().Add(time.Millisecond)
			timePerQuestionSec   = 8
		)
		cell := New(isCorrect, isFirstCorrectAnswer, qSentAt, answeredAt, timePerQuestionSec)
		require.Equal(t, 1, cell.FindSegmentNum())
		require.Equal(t, uint32(6), cell.PTS())
	}

	{
		var (
			isCorrect            = true
			isFirstCorrectAnswer = false
			qSentAt              = time.Now()
			answeredAt           = time.Now().Add(7 * time.Second)
			timePerQuestionSec   = 8
		)
		cell := New(isCorrect, isFirstCorrectAnswer, qSentAt, answeredAt, timePerQuestionSec)
		require.Equal(t, 4, cell.FindSegmentNum())
		require.Equal(t, uint32(1), cell.PTS())
	}

	{
		var (
			isCorrect            = false
			isFirstCorrectAnswer = false
			qSentAt              = time.Now()
			answeredAt           = time.Now().Add(time.Second)
			timePerQuestionSec   = 8
		)
		cell := New(isCorrect, isFirstCorrectAnswer, qSentAt, answeredAt, timePerQuestionSec)
		require.Equal(t, 1, cell.FindSegmentNum())
		require.Equal(t, uint32(0), cell.PTS())
	}
}
