CREATE TABLE votes (
    id         bigserial PRIMARY KEY,
    poll_id    INTEGER NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
    option     TEXT NOT NULL,
    created_at timestamp(0) WITH time zone NOT NULL DEFAULT now(),
    UNIQUE (poll_id)
);

CREATE INDEX idx_votes_poll_id ON votes(poll_id);
CREATE INDEX idx_votes_option ON votes(option);