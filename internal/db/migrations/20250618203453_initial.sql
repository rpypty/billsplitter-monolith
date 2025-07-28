-- +goose Up
-- +goose StatementBegin
-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id         INT       PRIMARY KEY NOT NULL,
    username   TEXT,
    first_name TEXT,
    last_name  TEXT,
    extra      JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX sessions_id ON sessions(id);
CREATE UNIQUE INDEX users_unique_extra_tg_id ON users ((extra->>'telegramID'));

-- Sessions table
CREATE TABLE sessions (
    id         UUID        PRIMARY KEY NOT NULL,
    user_id    INT         NOT NULL,
    expire_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX sessions_id ON sessions(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE sessions;
-- +goose StatementEnd
