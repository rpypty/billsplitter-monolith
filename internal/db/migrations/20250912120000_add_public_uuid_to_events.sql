-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN public_uuid UUID NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX events_public_uuid_uq ON events(public_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS events_public_uuid_uq;

ALTER TABLE events
    DROP COLUMN public_uuid;
-- +goose StatementEnd
