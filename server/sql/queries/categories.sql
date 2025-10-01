-- name: CreateCategory :one
INSERT INTO categories (user_id, name, color, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: GetCategoriesByUser :many
SELECT * FROM categories 
WHERE user_id = $1
ORDER BY name;

-- name: GetCategoryByID :one
SELECT * FROM categories 
WHERE id = $1 AND user_id = $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, color = $3
WHERE id = $1 AND user_id = $4
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories 
WHERE id = $1 AND user_id = $2;