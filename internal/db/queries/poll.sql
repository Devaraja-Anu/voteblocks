
-- name: CreatePoll :one
INSERT INTO polls (title, description, options, created_by,expires_at)
VALUES ($1, $2, $3, $4,$5)
RETURNING *;

-- name: GetPoll :one
SELECT * FROM polls WHERE id = $1;

-- name: ListPolls :many
SELECT * FROM polls
WHERE active = true AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetPollsByUser :many
SELECT * FROM polls
WHERE created_by = $1 
AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') ORDER BY created_at DESC;


