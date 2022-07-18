-- name: LinkDeviceToUser :exec
INSERT INTO users_devices (user_id, device_id) 
VALUES (@user_id, @device_id) 
ON CONFLICT (user_id, device_id) DO NOTHING;

-- name: GetUserIDsOnTheSameDevice :many 
SELECT user_id FROM users_devices
WHERE device_id IN (
    SELECT device_id 
    FROM users_devices 
    GROUP BY device_id
    HAVING count(user_id) > 1 
);

-- name: BlockUsersOnTheSameDevice :exec 
UPDATE users SET disabled = TRUE, block_reason = 'detected scam: multiple accounts on the same device'
WHERE id IN (
    SELECT user_id FROM users_devices
    WHERE device_id IN (
        SELECT users_devices.device_id 
        FROM users_devices 
        WHERE users_devices.device_id != @exclude_device_id
        GROUP BY users_devices.device_id
        HAVING count(users_devices.user_id) > 1 
    ) AND user_id IN (
        SELECT user_id FROM users WHERE verified_at IS NOT NULL
    )
) 
AND email NOT IN (SELECT allowed_value FROM whitelist WHERE allowed_type = 'email')
AND disabled = FALSE;

-- name: DoesUserHaveMoreThanOneAccount :one
SELECT count(t2.device_id) > 1 FROM users_devices AS t2
WHERE t2.device_id IN (
    SELECT t.device_id 
    FROM users_devices AS t
    WHERE t.user_id = @user_id
) GROUP BY t2.device_id;