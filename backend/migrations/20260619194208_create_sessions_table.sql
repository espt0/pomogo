-- +goose Up
CREATE TABLE sessions (
    id                UUID PRIMARY KEY,
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id           UUID NULL REFERENCES tasks(id) ON DELETE SET NULL,
    type              VARCHAR(50) NOT NULL
        CHECK (type IN ('work', 'short_break', 'long_break')),
    duration_minutes  INTEGER NOT NULL DEFAULT 0,
    started_at        TIMESTAMPTZ NOT NULL,
    ended_at          TIMESTAMPTZ NULL,
    completed         BOOLEAN NOT NULL DEFAULT false,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_session_completion
        CHECK (
            (completed = false AND ended_at IS NULL) OR
            (completed = true AND ended_at IS NOT NULL)
        )
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_task_id ON sessions(task_id);

-- +goose Down
DROP TABLE IF EXISTS sessions;
