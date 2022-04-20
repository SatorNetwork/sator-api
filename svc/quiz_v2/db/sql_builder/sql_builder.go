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
),
user_attempts_per_challenge AS (
	SELECT challenge_id, COUNT(*) as attempts FROM passed_challenges_data WHERE user_id='%v' GROUP BY challenge_id
),
passed_challenges_data AS (
	SELECT DISTINCT * FROM passed_challenges_data
)
SELECT
	challenges.id AS id,
	CONCAT(shows.title, '. ', episodes.title) AS title,
	challenges.players_to_start AS players_to_start,
	sorted_challenges.players_num AS players_num,
	challenges.prize_pool AS prize_pool,
	COALESCE((episode_access.activated_before > NOW())::BOOLEAN, FALSE) AS is_activated,
	episodes.cover AS cover,
	user_attempts_per_challenge.attempts AS attempts,
	challenges.user_max_attempts AS user_max_attempts
FROM
	challenges
	JOIN sorted_challenges ON sorted_challenges.id = challenges.id
		AND challenges.kind = 0
	LEFT JOIN passed_challenges_data ON passed_challenges_data.challenge_id = challenges.id
		AND passed_challenges_data.user_id = '%v'
	JOIN episodes ON episodes.challenge_id = challenges.id
	LEFT JOIN episode_access ON episode_access.episode_id = episodes.id
		AND episode_access.user_id = '%v'
		AND episode_access.activated_before > NOW()
	JOIN shows ON episodes.show_id = shows.id
	LEFT JOIN user_attempts_per_challenge ON user_attempts_per_challenge.challenge_id = challenges.id
WHERE (passed_challenges_data.reward_amount = 0
	OR passed_challenges_data.reward_amount IS NULL)
AND episodes.archived = FALSE
AND shows.archived = FALSE
AND (attempts IS NULL OR user_max_attempts - attempts > 0)
ORDER BY
	(challenges.players_to_start - sorted_challenges.players_num) ASC,
	sorted_challenges.players_num DESC,
	challenges.updated_at DESC
LIMIT %v OFFSET %v;
`

func ConstructGetChallengesSortedByPlayersQuery(userID uuid.UUID, limit, offset int32) string {
	return fmt.Sprintf(queryTemplate, userID, userID, userID, limit, offset)
}
