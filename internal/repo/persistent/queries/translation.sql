-- name: ReadHistory :many
SELECT source, destination, original, translation FROM history;

-- name: CreateHistory :exec
INSERT INTO history (source, destination, original, translation)
VALUES ($1, $2, $3, $4);
