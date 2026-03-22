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

-- name: UpdateTodo :one
UPDATE todos
SET title = $2, description = $3, completed = $4, priority = $5, due_date = $6, updated_at = NOW()
WHERE id = $1 AND user_id = $7 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTodo :exec
UPDATE todos
SET deleted_at = NOW()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;
