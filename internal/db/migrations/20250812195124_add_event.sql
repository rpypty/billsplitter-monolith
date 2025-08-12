-- +goose Up
-- +goose StatementBegin
CREATE TABLE event (
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

CREATE TABLE member (
    id          SERIAL PRIMARY KEY,
    user_id     INT,
    event_id    INT,
    member_name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event;
DROP TABLE member;
-- +goose StatementEnd
