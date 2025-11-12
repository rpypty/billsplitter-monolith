-- +goose Up
-- +goose StatementBegin
CREATE TABLE bills (
    id                 SERIAL        PRIMARY KEY,
    event_id           INT           NOT NULL REFERENCES events(id),
    name               TEXT          NOT NULL,
    created_by_user_id INT           NOT NULL REFERENCES users(id),
    total_amount       BIGINT        NOT NULL,
    currency           TEXT          NOT NULL,
    split_type         TEXT          NOT NULL,
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at         TIMESTAMPTZ
);

CREATE INDEX bills_event_id_idx ON bills(event_id);

CREATE TABLE bill_participants (
    id         SERIAL PRIMARY KEY,
    bill_id    INT    NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    member_id  INT    NOT NULL REFERENCES members(id),
    amount     BIGINT NOT NULL
);

CREATE INDEX bill_participants_bill_id_idx ON bill_participants(bill_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bill_participants;
DROP TABLE bills;
-- +goose StatementEnd
