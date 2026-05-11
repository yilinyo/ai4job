CREATE TABLE job_postings (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID         NOT NULL REFERENCES users(id),
    company_name     VARCHAR(120) NOT NULL,
    company_intro    TEXT,
    job_title        VARCHAR(120) NOT NULL,
    job_requirements TEXT,
    base             VARCHAR(120),
    salary_range     VARCHAR(120),
    source_platform  VARCHAR(120),
    source_url       VARCHAR(500),
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending',
    remark           TEXT,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_job_postings_user_id ON job_postings(user_id);
CREATE INDEX idx_job_postings_user_status ON job_postings(user_id, status);
CREATE INDEX idx_job_postings_user_source_platform ON job_postings(user_id, source_platform);
