-- +goose Up
-- +goose StatementBegin
CREATE INDEX user_id_idx ON store (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX user_id_idx ON store;
-- +goose StatementEnd
