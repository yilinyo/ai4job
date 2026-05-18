CREATE TABLE activity_logs (
    id           UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      VARCHAR(36)  NOT NULL REFERENCES users(id),
    actor_type   VARCHAR(40)  NOT NULL DEFAULT 'user',
    action       VARCHAR(40)  NOT NULL,
    target_type  VARCHAR(40)  NOT NULL,
    target_id    UUID         NOT NULL,
    target_label VARCHAR(200) NOT NULL,
    summary      TEXT         NOT NULL,
    route        VARCHAR(500) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_logs_user_created_at ON activity_logs(user_id, created_at DESC);
CREATE INDEX idx_activity_logs_target ON activity_logs(target_type, target_id);
