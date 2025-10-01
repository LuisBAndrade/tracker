-- name: CreateExpense :one
INSERT INTO expenses (user_id, category_id, amount, description, date, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING *;

-- name: GetExpensesByUser :many
SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.date, e.created_at, e.updated_at,
       c.name as category_name, c.color as category_color
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1
ORDER BY e.date DESC, e.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetExpensesByUserAndDateRange :many
SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.date, e.created_at, e.updated_at,
       c.name as category_name, c.color as category_color
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.user_id = $1 AND e.date BETWEEN $2 AND $3
ORDER BY e.date DESC, e.created_at DESC;

-- name: GetExpenseByID :one
SELECT e.id, e.user_id, e.category_id, e.amount, e.description, e.date, e.created_at, e.updated_at,
       c.name as category_name, c.color as category_color
FROM expenses e
LEFT JOIN categories c ON e.category_id = c.id
WHERE e.id = $1 AND e.user_id = $2;

-- name: UpdateExpense :one
UPDATE expenses
SET amount = $3, description = $4, category_id = $5, date = $6, updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1 AND user_id = $2;

-- name: GetExpenseTotalByUser :one
SELECT COALESCE(SUM(amount), 0) as total
FROM expenses
WHERE user_id = $1;

-- name: GetExpenseTotalByUserAndDateRange :one
SELECT COALESCE(SUM(amount), 0) as total
FROM expenses
WHERE user_id = $1 AND date BETWEEN $2 AND $3;

-- name: GetExpensesByCategory :many
SELECT 
    c.id as category_id,
    c.name as category_name,
    c.color as category_color,
    COALESCE(SUM(e.amount), 0) as total_amount,
    COUNT(e.id) as expense_count
FROM categories c
LEFT JOIN expenses e ON c.id = e.category_id AND e.user_id = $1 AND e.date BETWEEN $2 AND $3
WHERE c.user_id = $1
GROUP BY c.id, c.name, c.color
HAVING COUNT(e.id) > 0 OR $4 = true
ORDER BY total_amount DESC;