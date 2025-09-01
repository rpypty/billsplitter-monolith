-- +goose Up
-- +goose StatementBegin

-- Payment methods table
CREATE TABLE IF NOT EXISTS payment_methods (
    id         SERIAL        PRIMARY KEY,
    user_id    INT           NOT NULL,
    name       TEXT          NOT NULL,
    recipient  TEXT          NOT NULL,
    description TEXT         NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX payment_methods_user_id ON payment_methods(user_id);
CREATE INDEX payment_methods_id ON payment_methods(id);

-- Foreign key constraint
ALTER TABLE payment_methods 
ADD CONSTRAINT fk_payment_methods_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payment_methods;
-- +goose StatementEnd 