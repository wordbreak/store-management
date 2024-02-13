-- +goose Up
-- +goose StatementBegin
-- add column abstract_name varchar(255) not null;
ALTER TABLE product ADD COLUMN abstract_name VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product DROP COLUMN abstract_name;
-- +goose StatementEnd
