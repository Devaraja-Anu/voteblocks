CREATE TABLE polls (
    id           bigserial PRIMARY KEY,
    title        text NOT NULL,
    description  TEXT NOT NULL,
    options      TEXT[] NOT NULL,
    -- created_by   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   timestamp(0) WITH time zone NOT NULL DEFAULT now(),
    expires_at   timestamp(0) WITH time zone DEFAULT NULL,
    active       BOOLEAN NOT NULL DEFAULT true
);

-- CREATE INDEX idx_polls_created_by ON polls(created_by);
CREATE INDEX idx_polls_expires_at ON polls(expires_at);
CREATE INDEX idx_polls_active ON polls(active);
