CREATE TABLE interview_progress (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID        NOT NULL REFERENCES users(id),
    application_id      UUID        NOT NULL REFERENCES resume_applications(id),
    stage               VARCHAR(30) NOT NULL,
    occurred_at         TIMESTAMPTZ NOT NULL,
    result              VARCHAR(20) NOT NULL,
    interviewer_or_team VARCHAR(120),
    note                TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_interview_progress_user_id ON interview_progress(user_id);
CREATE INDEX idx_interview_progress_application_id ON interview_progress(application_id);
CREATE INDEX idx_interview_progress_occurred_at ON interview_progress(occurred_at);
