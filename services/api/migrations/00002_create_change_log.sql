-- +goose Up
CREATE TABLE user_change_counters (
    user_id TEXT PRIMARY KEY,
    last_sequence BIGINT NOT NULL DEFAULT 0 CHECK (last_sequence >=0)
);

CREATE TABLE change_log (
    user_id TEXT NOT NULL REFERENCES user_change_counters(user_id),
    sequence BIGINT NOT NULL CHECK (sequence > 0),
    entity_type TEXT NOT NULL,
    entity_id UUID NOT NULL,
    operation TEXT NOT NULL CHECK (operation IN ('create', 'update', 'delete')),
    revision BIGINT NOT NULL CHECK (revision > 0),
    payload JSONB NOT NULL CHECK (jsonb_typeof(payload) = 'object'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, sequence)

);

-- +goose Down
DROP TABLE change_log;
DROP TABLE user_change_counters;
