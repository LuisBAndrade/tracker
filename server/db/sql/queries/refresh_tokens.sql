-- name: StoreRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (?, ?, ?);

-- name: GetRefreshTokens :one
SELECT token, user_id, created_at, expires_at
FROM refresh_tokens
WHERE token = ? LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token = ?;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens WHERE user_id = ?;