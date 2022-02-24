package sql_builder

import "fmt"

const queryTemplate = `
SELECT id, COUNT(user_id) AS players_num FROM challenges
LEFT JOIN room_players ON room_players.challenge_id = challenges.id
GROUP BY id
ORDER BY players_num DESC
LIMIT %v OFFSET %v;
`

func ConstructGetChallengesSortedByPlayersQuery(limit, offset int32) string {
	return fmt.Sprintf(queryTemplate, limit, offset)
}
