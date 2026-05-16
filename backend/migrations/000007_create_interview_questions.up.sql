CREATE TABLE interview_questions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          VARCHAR(36) NOT NULL REFERENCES users(id),
    title            VARCHAR(200) NOT NULL,
    content          TEXT NOT NULL,
    answer_analysis  TEXT NOT NULL,
    tags             JSONB NOT NULL DEFAULT '[]',
    difficulty       VARCHAR(20) NOT NULL,
    question_type    VARCHAR(30) NOT NULL,
    mastery_status   VARCHAR(30) NOT NULL DEFAULT 'unreviewed',
    is_favorite      BOOLEAN NOT NULL DEFAULT FALSE,
    source           VARCHAR(100),
    last_reviewed_at TIMESTAMPTZ,
    remark           TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_interview_questions_user_id ON interview_questions(user_id);
CREATE INDEX idx_interview_questions_user_difficulty ON interview_questions(user_id, difficulty) WHERE deleted_at IS NULL;
CREATE INDEX idx_interview_questions_user_question_type ON interview_questions(user_id, question_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_interview_questions_user_mastery_status ON interview_questions(user_id, mastery_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_interview_questions_user_is_favorite ON interview_questions(user_id, is_favorite) WHERE deleted_at IS NULL;
CREATE INDEX idx_interview_questions_user_created_at ON interview_questions(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_interview_questions_tags_gin ON interview_questions USING GIN(tags);

CREATE TABLE question_occurrences (
    id                       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id                  VARCHAR(36) NOT NULL REFERENCES users(id),
    question_id              UUID NOT NULL REFERENCES interview_questions(id),
    application_id           UUID NOT NULL REFERENCES resume_applications(id),
    interview_progress_id    UUID REFERENCES interview_progress(id),
    company_name_snapshot    VARCHAR(200) NOT NULL,
    job_title_snapshot       VARCHAR(200) NOT NULL,
    interview_stage_snapshot VARCHAR(50),
    occurred_at              TIMESTAMPTZ NOT NULL,
    actual_question          TEXT,
    answer_performance       VARCHAR(20),
    note                     TEXT,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at               TIMESTAMPTZ
);

CREATE INDEX idx_question_occurrences_user_id ON question_occurrences(user_id);
CREATE INDEX idx_question_occurrences_question_id ON question_occurrences(question_id);
CREATE INDEX idx_question_occurrences_application_id ON question_occurrences(application_id);
CREATE INDEX idx_question_occurrences_interview_progress_id ON question_occurrences(interview_progress_id);
CREATE INDEX idx_question_occurrences_occurred_at ON question_occurrences(occurred_at DESC);
CREATE UNIQUE INDEX idx_question_occurrences_user_question_application_progress_unique
ON question_occurrences(user_id, question_id, application_id, interview_progress_id)
WHERE interview_progress_id IS NOT NULL AND deleted_at IS NULL;

CREATE TABLE flashcard_reviews (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     VARCHAR(36) NOT NULL REFERENCES users(id),
    question_id UUID NOT NULL REFERENCES interview_questions(id),
    reviewed_at TIMESTAMPTZ NOT NULL,
    result      VARCHAR(20) NOT NULL,
    review_mode VARCHAR(30) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_flashcard_reviews_user_id ON flashcard_reviews(user_id);
CREATE INDEX idx_flashcard_reviews_question_id ON flashcard_reviews(question_id);
CREATE INDEX idx_flashcard_reviews_reviewed_at ON flashcard_reviews(reviewed_at DESC);
