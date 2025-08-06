
-- name: CreatePoll :one
INSERT INTO polls (title,description,options,expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

DELETE FROM polls WHERE id = $1;

-- name: GetPoll :one
SELECT * FROM polls WHERE id = $1;

-- name: ListPolls :many
SELECT * FROM polls
WHERE active = true AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC LIMIT $1 OFFSET $2;


-- name: DeactivateExpiredPolls :exec
UPDATE polls
SET active = false
WHERE active = true
  AND expires_at IS NOT NULL
  AND expires_at < now();
