-- name: CreateAnnouncement :one
INSERT INTO announcements (
    title,
    description,
    action_url,
    starts_at,
    ends_at,
    type,
    type_specific_params
)
VALUES (
    @title,
    @description,
    @action_url,
    @starts_at,
    @ends_at,
    @type,
    @type_specific_params
) RETURNING *;

-- name: GetAnnouncementByID :one
SELECT * FROM announcements
WHERE id = @id;

-- name: UpdateAnnouncementByID :exec
UPDATE announcements
SET
    title = @title,
    description = @description,
    action_url = @action_url,
    starts_at = @starts_at,
    ends_at = @ends_at,
    type = @type,
    type_specific_params = @type_specific_params
WHERE id = @id;

-- name: DeleteAnnouncementByID :exec
DELETE FROM announcements
WHERE id = @id;

-- name: ListAnnouncements :many
SELECT * FROM announcements;

-- name: ListUnreadAnnouncements :many
SELECT * FROM announcements WHERE id IN (
    SELECT id
    FROM announcements
        EXCEPT
    SELECT announcement_id
    FROM read_announcements
    WHERE user_id = @user_id
);

-- name: ListActiveAnnouncements :many
SELECT * FROM announcements
WHERE starts_at <= NOW() AND NOW() <= ends_at;

-- name: CleanUpReadAnnouncements :exec
DELETE FROM read_announcements;

-- name: CleanUpAnnouncements :exec
DELETE FROM announcements;
