-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
    id                 SERIAL        PRIMARY KEY,
    name               TEXT,
    created_by_user_id INT,
    status             TEXT,
    event_date         TIMESTAMPTZ,
    event_type         TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ
);

CREATE TABLE members (
    id          SERIAL PRIMARY KEY,
    user_id     INT,
    event_id    INT,
    member_name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
DROP TABLE membesr;
-- +goose StatementEnd
