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
UPDATE users SET disabled = TRUE, block_reason = 'suspicion of fraud: created multiple accounts'
WHERE device_id IN (
    SELECT device_id 
    FROM users_devices 
    GROUP BY device_id
    HAVING count(user_id) > 1 
);