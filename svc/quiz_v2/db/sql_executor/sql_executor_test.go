package sql_executor

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/db/sql_builder"
)

func TestGetChallengesSortedByPlayers(t *testing.T) {
	dbConnString := "postgresql://pguser:pgpass@127.0.0.1/pgdb?sslmode=disable"
	dbClient, err := sql.Open("postgres", dbConnString)
	require.NoError(t, err)

	query := sql_builder.ConstructGetChallengesSortedByPlayersQuery(1, 0)
	sqlExecutor := New(dbClient)
	_, err = sqlExecutor.ExecuteGetChallengesSortedByPlayersQuery(query, nil)
	require.NoError(t, err)
}
