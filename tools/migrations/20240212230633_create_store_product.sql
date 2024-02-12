-- +goose Up
-- +goose StatementBegin
CREATE TABLE store_product (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    store_id BIGINT NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store_product;
-- +goose StatementEnd
