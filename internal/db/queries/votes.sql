-- name: CreateVote :one
INSERT INTO votes (user_id, poll_id, option)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetVotesByPoll :many
SELECT * FROM votes WHERE poll_id = $1;

-- name: GetVoteByUserAndPoll :one
SELECT * FROM votes WHERE user_id = $1 AND poll_id = $2;

-- name: CountVotesByOption :many
SELECT option, COUNT(*) as vote_count
FROM votes
WHERE poll_id = $1
GROUP BY option
ORDER BY vote_count DESC;
