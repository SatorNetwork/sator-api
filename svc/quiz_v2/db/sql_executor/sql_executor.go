package sql_executor

import (
	"database/sql"

	"github.com/pkg/errors"
)

type SQLExecutor struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLExecutor {
	return &SQLExecutor{
		db: db,
	}
}

type Challenge struct {
	ID             string
	Title          string
	PlayersToStart int
	PlayersNum     int
	PrizePool      float64
	IsActivated    bool
	Cover          string
}

func (e *SQLExecutor) ExecuteGetChallengesSortedByPlayersQuery(sql string, args []interface{}) ([]*Challenge, error) {
	challenges := make([]*Challenge, 0)
	rows, err := e.db.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "can't execute query")
	}
	defer rows.Close()
	for rows.Next() {
		challenge := new(Challenge)
		err := rows.Scan(
			&challenge.ID,
			&challenge.Title,
			&challenge.PlayersToStart,
			&challenge.PlayersNum,
			&challenge.PrizePool,
			&challenge.IsActivated,
			&challenge.Cover,
		)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, challenge)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}
