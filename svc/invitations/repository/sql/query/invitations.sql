-- name: GetInvitationByInviteeEmail :one
SELECT *
FROM invitations
WHERE normalized_invitee_email = $1
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
-- name: GetInvitationsByInvitedByID :many
SELECT *
FROM invitations
WHERE invited_by = $1
ORDER BY invited_at DESC;
-- name: CreateInvitation :one
INSERT INTO invitations (invitee_email, normalized_invitee_email, invited_by, accepted_by)
VALUES ($1, $2, $3, uuid_nil()) RETURNING *;
-- name: AcceptInvitationByInviteeEmail :exec
UPDATE invitations
SET accepted_by = @accepted_by,
    accepted_at = @accepted_at
WHERE invitee_email = @invitee_email;
