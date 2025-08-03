CREATE TABLE polls (
    id           SERIAL PRIMARY KEY,
    title        TEXT NOT NULL,
    description  TEXT,
    options      TEXT[] NOT NULL,
    created_by   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    expires_at   TIMESTAMP,
    active       BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_polls_created_by ON polls(created_by);
CREATE INDEX idx_polls_expires_at ON polls(expires_at);
CREATE INDEX idx_polls_active ON polls(active);
