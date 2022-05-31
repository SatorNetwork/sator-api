-- name: DisableNotification :exec
INSERT INTO disabled_notifications (
    user_id,
    topic
)
VALUES (
    @user_id,
    @topic
);

-- name: EnableNotification :exec
DELETE FROM disabled_notifications
WHERE user_id = @user_id AND topic = @topic;

-- name: IsNotificationEnabled :one
SELECT (NOT EXISTS(
    SELECT * FROM disabled_notifications
    WHERE user_id = @user_id AND topic = @topic
))::BOOL;

-- name: IsNotificationDisabled :one
SELECT EXISTS(
    SELECT * FROM disabled_notifications
    WHERE user_id = @user_id AND topic = @topic
);