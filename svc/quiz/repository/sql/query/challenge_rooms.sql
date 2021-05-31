-- name: AddNewChallegeRoom :one
INSERT INTO challenge_rooms (
        challenge_id,
        prize_pool,
        players_to_start,
        time_per_question
    )
VALUES (
        @challenge_id,
        @prize_pool,
        @players_to_start,
        @time_per_question
    ) RETURNING *;
-- name: UpdateChallengeRoomStatus :exec
UPDATE challenge_rooms
SET status = @status
WHERE id = @id;
-- name: GetChallengeRoomByID :one
SELECT *
FROM challenge_rooms
WHERE id = $1;
-- name: GetChallengeRoomByChallengeID :one
SELECT *
FROM challenge_rooms
WHERE challenge_id = $1
    AND status = 0
ORDER BY created_at DESC;