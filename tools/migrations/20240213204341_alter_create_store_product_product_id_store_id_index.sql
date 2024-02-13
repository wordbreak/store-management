-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_product_id_store_id_idx ON store_product (product_id, store_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uniq_product_id_store_id_idx ON store_product;
-- +goose StatementEnd
