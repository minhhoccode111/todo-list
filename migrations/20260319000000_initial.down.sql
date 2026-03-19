DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_deleted_at;
DROP INDEX IF EXISTS idx_todos_due_date;
DROP INDEX IF EXISTS idx_todos_completed;
DROP INDEX IF EXISTS idx_todos_user_id;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS priority_level;
