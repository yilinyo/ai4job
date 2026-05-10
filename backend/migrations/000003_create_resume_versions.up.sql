CREATE TABLE resume_versions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    repo_id     UUID         NOT NULL REFERENCES resume_repos(id),
    title       VARCHAR(120) NOT NULL,
    content     TEXT         NOT NULL,
    version_num INT          NOT NULL,
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resume_versions_repo_version_unique
ON resume_versions(repo_id, version_num)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_resume_versions_repo_default_unique
ON resume_versions(repo_id)
WHERE is_default = TRUE AND deleted_at IS NULL;

CREATE INDEX idx_resume_versions_repo_id ON resume_versions(repo_id);
