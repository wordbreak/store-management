-- +goose Up
-- +goose StatementBegin
CREATE TABLE store (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL
) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE store;
-- +goose StatementEnd
