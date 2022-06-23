-- name: MarkAsRead :exec
INSERT INTO read_announcements (
    announcement_id,
    user_id
)
VALUES (
    @announcement_id,
    @user_id
) RETURNING *;

-- name: IsRead :one
SELECT EXISTS(
    SELECT * FROM read_announcements
    WHERE announcement_id = @announcement_id AND
          user_id = @user_id
);

-- name: IsNotRead :one
SELECT (NOT EXISTS(
   SELECT * FROM read_announcements
   WHERE announcement_id = @announcement_id AND
         user_id = @user_id
))::BOOL;
