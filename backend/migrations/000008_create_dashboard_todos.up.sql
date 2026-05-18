CREATE TABLE dashboard_todos (
    id           UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      VARCHAR(36)  NOT NULL REFERENCES users(id),
    title        VARCHAR(200) NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    due_at       TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    target_type  VARCHAR(40)  NOT NULL DEFAULT '',
    target_id    UUID,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dashboard_todos_user_id ON dashboard_todos(user_id);
CREATE INDEX idx_dashboard_todos_user_completed ON dashboard_todos(user_id, completed_at);
CREATE INDEX idx_dashboard_todos_user_due_at ON dashboard_todos(user_id, due_at);
CREATE INDEX idx_dashboard_todos_target ON dashboard_todos(target_type, target_id);
