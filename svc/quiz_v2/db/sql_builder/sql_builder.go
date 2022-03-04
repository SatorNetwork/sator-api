package sql_builder

import "fmt"

const queryTemplate = `
WITH sorted_challenges AS (
	SELECT
		id,
		COUNT(user_id) AS players_num
	FROM
		challenges
	LEFT JOIN room_players ON room_players.challenge_id = challenges.id
GROUP BY
	id
ORDER BY
	players_num DESC
)
SELECT
	challenges.id as id,
	challenges.title as title,
	challenges.players_to_start as players_to_start,
	sorted_challenges.players_num AS players_num,
	challenges.prize_pool as prize_pool
FROM
	challenges
	JOIN sorted_challenges ON sorted_challenges.id = challenges.id
		AND challenges.kind = 0
LIMIT %v OFFSET %v;
`

func ConstructGetChallengesSortedByPlayersQuery(limit, offset int32) string {
	return fmt.Sprintf(queryTemplate, limit, offset)
}
