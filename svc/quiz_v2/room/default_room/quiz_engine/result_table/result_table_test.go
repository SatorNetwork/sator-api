package result_table

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestResultTable(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	if false {
		cfg := Config{
			QuestionNum:        5,
			WinnersNum:         2,
			PrizePool:          250,
			TimePerQuestionSec: 8,
		}
		rt := New(&cfg)

		for qNum := 0; qNum < cfg.QuestionNum; qNum++ {
			err := rt.RegisterQuestionSendingEvent(qNum)
			require.NoError(t, err)
			if qNum == 0 {
				err = rt.RegisterAnswer(userID1, qNum, true, time.Now().Add(time.Millisecond))
				require.NoError(t, err)
				err = rt.RegisterAnswer(userID2, qNum, true, time.Now().Add(time.Millisecond))
				require.NoError(t, err)
			}
		}

		userIDToPrice := rt.GetPrizePoolDistribution()
		require.NotNil(t, userIDToPrice)
	}

	{
		cfg := Config{
			QuestionNum:        5,
			WinnersNum:         2,
			PrizePool:          250,
			TimePerQuestionSec: 8,
		}
		rt := New(&cfg)

		for qNum := 0; qNum < cfg.QuestionNum; qNum++ {
			err := rt.RegisterQuestionSendingEvent(qNum)
			require.NoError(t, err)
			err = rt.RegisterAnswer(userID1, qNum, true, time.Now().Add(time.Millisecond))
			require.NoError(t, err)
			err = rt.RegisterAnswer(userID2, qNum, true, time.Now().Add(time.Millisecond))
			require.NoError(t, err)
		}

		ptsMap := map[uuid.UUID]uint32{
			userID1: 30,
			userID2: 20,
		}
		require.Equal(t, ptsMap, rt.calcPTSMap())
		require.Equal(t, []uuid.UUID{userID1, userID2}, rt.getWinnerIDs())
		require.Equal(t, ptsMap, rt.calcWinnersMap())

		userIDToPrice := map[uuid.UUID]float64{
			userID1: 150,
			userID2: 100,
		}
		require.Equal(t, userIDToPrice, rt.GetPrizePoolDistribution())
	}

	{
		cfg := Config{
			QuestionNum:        5,
			WinnersNum:         2,
			PrizePool:          250,
			TimePerQuestionSec: 8,
		}
		rt := New(&cfg)

		for qNum := 0; qNum < cfg.QuestionNum; qNum++ {
			err := rt.RegisterQuestionSendingEvent(qNum)
			require.NoError(t, err)
			err = rt.RegisterAnswer(userID1, qNum, true, time.Now().Add(time.Millisecond))
			require.NoError(t, err)
			err = rt.RegisterAnswer(userID2, qNum, true, time.Now().Add(time.Millisecond))
			require.NoError(t, err)
			err = rt.RegisterAnswer(userID3, qNum, true, time.Now().Add(7*time.Second))
			require.NoError(t, err)
		}

		ptsMap := map[uuid.UUID]uint32{
			userID1: 30,
			userID2: 20,
			userID3: 5,
		}
		require.Equal(t, ptsMap, rt.calcPTSMap())
		require.Equal(t, []uuid.UUID{userID1, userID2}, rt.getWinnerIDs())
		winnersMap := map[uuid.UUID]uint32{
			userID1: 30,
			userID2: 20,
		}
		require.Equal(t, winnersMap, rt.calcWinnersMap())

		userIDToPrice := map[uuid.UUID]float64{
			userID1: 150,
			userID2: 100,
		}
		require.Equal(t, userIDToPrice, rt.GetPrizePoolDistribution())
	}
}
