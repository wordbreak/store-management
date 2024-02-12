-- +goose Up
-- +goose StatementBegin
CREATE TABLE product (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(255) NOT NULL,
    price DECIMAL(20,4) UNSIGNED NOT NULL,
    cost DECIMAL(20,4) UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    barcode VARCHAR(255) NOT NULL,
    expiry_date DATE NOT NULL,
    size VARCHAR(64) NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product;
-- +goose StatementEnd
