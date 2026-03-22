-- name: ReadTodos :many
SELECT id, user_id, title, description, completed, priority, due_date, created_at, updated_at
FROM todos t
WHERE deleted_at IS NOT NULL
AND t.user_id = $1
LIMIT $2
OFFSET $3;

-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description, priority, due_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
