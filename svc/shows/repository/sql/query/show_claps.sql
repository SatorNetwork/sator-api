-- name: CountUserClaps :one
SELECT COUNT(*) FROM show_claps
WHERE show_id = @show_id
AND user_id = @user_id
GROUP BY show_id;

-- name: AddClapForShow :exec
INSERT INTO show_claps (show_id, user_id)
VALUES (@show_id, @user_id);