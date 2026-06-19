-- +goose Up
CREATE TABLE settings (
    id                     UUID PRIMARY KEY,
    user_id                UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    work_duration          INTEGER NOT NULL DEFAULT 25,
    short_break_duration   INTEGER NOT NULL DEFAULT 5,
    long_break_duration    INTEGER NOT NULL DEFAULT 15,
    long_break_interval    INTEGER NOT NULL DEFAULT 4,
    auto_start_work        BOOLEAN NOT NULL DEFAULT false,
    auto_start_break       BOOLEAN NOT NULL DEFAULT false,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS settings;
