package sql_builder

import (
	"fmt"

	"github.com/google/uuid"
)

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
	challenges.id AS id,
	challenges.title AS title,
	challenges.players_to_start AS players_to_start,
	sorted_challenges.players_num AS players_num,
	challenges.prize_pool AS prize_pool
FROM
	challenges
	JOIN sorted_challenges ON sorted_challenges.id = challenges.id
		AND challenges.kind = 0
	LEFT JOIN passed_challenges_data ON passed_challenges_data.challenge_id = challenges.id
		AND passed_challenges_data.user_id = '%v'
WHERE
	passed_challenges_data.reward_amount = 0
	OR passed_challenges_data.reward_amount IS NULL
ORDER BY
	(challenges.players_to_start - sorted_challenges.players_num) ASC,
	sorted_challenges.players_num DESC,
	challenges.updated_at DESC
LIMIT %v OFFSET %v;
`

func ConstructGetChallengesSortedByPlayersQuery(userID uuid.UUID, limit, offset int32) string {
	return fmt.Sprintf(queryTemplate, userID, limit, offset)
}
