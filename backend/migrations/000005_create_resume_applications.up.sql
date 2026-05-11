CREATE TABLE resume_applications (
    id                          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                     VARCHAR(36)  NOT NULL REFERENCES users(id),
    job_posting_id              UUID         REFERENCES job_postings(id),
    resume_version_id           UUID         NOT NULL REFERENCES resume_versions(id),
    repo_id                     UUID         NOT NULL REFERENCES resume_repos(id),
    company_name                VARCHAR(120) NOT NULL,
    company_intro               TEXT,
    job_title                   VARCHAR(120) NOT NULL,
    job_requirements            TEXT,
    base                        VARCHAR(120),
    salary_range                VARCHAR(120),
    source_platform             VARCHAR(120),
    source_url                  VARCHAR(500),
    delivery_channel            VARCHAR(120),
    applied_at                  TIMESTAMPTZ  NOT NULL,
    current_status              VARCHAR(30)  NOT NULL DEFAULT 'applied',
    interview_summary           TEXT,
    interview_notes             TEXT,
    resume_match_score          INT          NOT NULL DEFAULT 0,
    interview_performance_score INT          NOT NULL DEFAULT 0,
    role_interest_score         INT          NOT NULL DEFAULT 0,
    overall_score               INT          NOT NULL DEFAULT 0,
    remark                      TEXT,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ
);

CREATE INDEX idx_resume_applications_user_id ON resume_applications(user_id);
CREATE INDEX idx_resume_applications_job_posting_id ON resume_applications(job_posting_id);
CREATE INDEX idx_resume_applications_resume_version_id ON resume_applications(resume_version_id);
CREATE INDEX idx_resume_applications_repo_id ON resume_applications(repo_id);
CREATE INDEX idx_resume_applications_current_status ON resume_applications(current_status);
CREATE INDEX idx_resume_applications_applied_at ON resume_applications(applied_at);
