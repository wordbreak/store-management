-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  phone_number VARCHAR(64) UNIQUE NOT NULL,
  password VARCHAR(512) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
