-- name: AddVote :one
INSERT INTO votes (poll_id,option) VALUES ($1,$2)
RETURNING id,poll_id,option,created_at;

-- name: GetTotalVotesForPoll :one
SELECT count(*) AS total FROM votes WHERE poll_id = $1;
