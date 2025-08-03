-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetPollsVotedByUser :many
SELECT
  p.id,
  p.title,
  p.description,
  p.options,
  p.created_by,
  p.created_at,
  p.expires_at,
  p.active,
  v.option AS voted_option
FROM polls p
JOIN votes v ON p.id = v.poll_id
WHERE v.user_id = $1;