-- name: GetInvitationByInviteeEmail :one
SELECT *
FROM invitations
WHERE email = $1
    LIMIT 1;
-- name: GetInvitationByInviteeID :one
SELECT *
FROM invitations
WHERE accepted_by = $1
    LIMIT 1;
-- name: GetInvitations :many
SELECT *
FROM invitations
ORDER BY invited_at DESC;
-- name: GetInvitationsPaginated :many
SELECT *
FROM invitations
ORDER BY invited_at DESC
LIMIT $1 OFFSET $2;
-- name: GetInvitationsByInviterID :many
SELECT *
FROM invitations
WHERE invited_by = $1
ORDER BY invited_at DESC;
-- name: CreateInvitation :one
INSERT INTO invitations (email, invited_by)
VALUES ($1, $2) RETURNING *;
-- name: AcceptInvitationByInviteeEmail :exec
UPDATE invitations
SET accepted_by = @accepted_by,
    accepted_at = @accepted_at
WHERE id = @id;
-- name: SetRewardReceived :exec
UPDATE invitations
SET reward_received = @reward_received
WHERE id = @id;
