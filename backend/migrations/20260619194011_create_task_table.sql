-- +goose Up
CREATE TABLE task (
    id                   UUID PRIMARY KEY,
    user_id              UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                VARCHAR(255) NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    status               VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'in_progress', 'completed', 'archived')),
    estimated_pomodoros  INTEGER NOT NULL DEFAULT 0,
    completed_pomodoros  INTEGER NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_task_user_id ON task(user_id);
CREATE INDEX idx_task_status ON task(status);

-- +goose Down
DROP TABLE IF EXISTS task;
