-- +goose Up
-- +goose StatementBegin
ALTER TABLE bills
    ADD COLUMN paid_by INT REFERENCES members(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bills
    DROP COLUMN paid_by;
-- +goose StatementEnd
