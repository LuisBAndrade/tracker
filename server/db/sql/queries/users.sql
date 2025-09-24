-- name: CreateUser :execresult
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    UUID(),
    NOW(),
    NOW(),
    ?,
    ?
);

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password FROM users
WHERE email = ? 
LIMIT 1;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE id = ?
LIMIT 1;