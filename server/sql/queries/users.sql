-- db/queries/auth.sql
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateSession :exec
INSERT INTO sessions (token, user_id, expires_at, created_at)
VALUES ($1, $2, $3, NOW());

-- name: GetUserBySessionToken :one
SELECT u.* FROM users u
JOIN sessions s ON u.id = s.user_id
WHERE s.token = $1 AND s.expires_at > NOW();

-- name: RevokeSession :exec
DELETE FROM sessions WHERE token = $1;

-- name: RevokeAllUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: CleanupExpiredSessions :exec
DELETE FROM sessions WHERE expires_at <= NOW();