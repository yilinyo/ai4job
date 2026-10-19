CREATE TABLE resume_attachments (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version_id  UUID         NOT NULL REFERENCES resume_versions(id) ON DELETE CASCADE,
    file_name   VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_size   BIGINT       NOT NULL CHECK (file_size >= 0),
    data        BYTEA        NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resume_attachments_version_id ON resume_attachments(version_id);
