
-- name: CreatePoll :one
INSERT INTO polls (title,description,options,expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

DELETE FROM polls WHERE id = $1;

-- name: GetPoll :one
SELECT * FROM polls WHERE id = $1;

-- name: ListPolls :many
SELECT count(*) OVER() AS total_records,
    id, title, description, options, created_at, expires_at, active FROM polls
WHERE active = true AND (expires_at IS NULL OR expires_at > now())
AND (to_tsvector('simple',title) @@ plainto_tsquery('simple',$1) OR $1 = '') 
ORDER BY created_at DESC LIMIT $2 OFFSET $3;


-- name: DeactivateExpiredPolls :exec
UPDATE polls
SET active = false
WHERE active = true
AND expires_at IS NOT NULL
AND expires_at < now();


-- name: GetPollWithVoteCounts :one
SELECT 
    p.id AS poll_id,
    p.title,
    p.description,
    p.options,
    p.created_at,
    p.expires_at,
    p.active,
    COALESCE(vote_counts.counts, '{}'::jsonb) AS vote_counts
FROM polls p
LEFT JOIN LATERAL (
    SELECT jsonb_object_agg(option, count) AS counts
    FROM (
        SELECT option, COUNT(*) AS count
        FROM votes
        WHERE poll_id = p.id
        GROUP BY option
    ) AS counted_votes
) AS vote_counts ON TRUE
WHERE p.id = $1;
