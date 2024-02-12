-- +goose Up
-- +goose StatementBegin
ALTER TABLE product MODIFY COLUMN expiry_date TIMESTAMP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product MODIFY COLUMN expiry_date DATE NOT NULL;
-- +goose StatementEnd
