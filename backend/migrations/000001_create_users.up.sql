CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id              VARCHAR(36)  PRIMARY KEY,
    username        VARCHAR(50)  NOT NULL,
    password_hash   VARCHAR      NOT NULL,
    email           VARCHAR(100) NOT NULL,
    real_name       VARCHAR(100) NOT NULL DEFAULT '',
    age             INT          NOT NULL DEFAULT 0,
    phone           VARCHAR(20)  NOT NULL DEFAULT '',
    job_intention   VARCHAR(255) NOT NULL DEFAULT '',
    avatar          VARCHAR(500) NOT NULL DEFAULT '',
    location        VARCHAR(255) NOT NULL DEFAULT '',
    personal_advantage TEXT      NOT NULL DEFAULT '',
    education_experiences JSONB  NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ  NULL
);

-- Unique index on username excluding soft-deleted rows
CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique
    ON users (username)
    WHERE deleted_at IS NULL;

-- Unique index on email excluding soft-deleted rows
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique
    ON users (email)
    WHERE deleted_at IS NULL;
