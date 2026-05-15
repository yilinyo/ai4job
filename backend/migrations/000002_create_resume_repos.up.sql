CREATE TABLE resume_repos (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         VARCHAR(36)  NOT NULL REFERENCES users(id),
    name            VARCHAR(100) NOT NULL,
    job_type        VARCHAR(20)  NOT NULL,
    target_position VARCHAR(100) NOT NULL,
    description     TEXT,
    version_count   INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_repos_user_name_unique
ON resume_repos(user_id, name)
WHERE deleted_at IS NULL;

CREATE INDEX idx_resume_repos_user_id ON resume_repos(user_id);
CREATE INDEX idx_resume_repos_user_job_type ON resume_repos(user_id, job_type);
