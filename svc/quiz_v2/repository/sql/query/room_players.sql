-- name: RegisterNewPlayer :exec
INSERT INTO room_players (
    challenge_id,
    user_id
)
VALUES (
    @challenge_id,
    @user_id
) ON CONFLICT (challenge_id, user_id) DO NOTHING;

-- name: UnregisterPlayer :exec
DELETE FROM room_players
WHERE challenge_id = @challenge_id
  AND user_id = @user_id;

-- name: CountPlayersInRoom :one
SELECT COUNT(user_id) AS players
FROM room_players
WHERE challenge_id = @challenge_id
GROUP BY challenge_id
LIMIT 1;
