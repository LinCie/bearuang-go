-- +goose Up
CREATE TABLE users (
    id VARCHAR(26) PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT user_email_key UNIQUE (email)
);

-- +goose Down
DROP TABLE users;
