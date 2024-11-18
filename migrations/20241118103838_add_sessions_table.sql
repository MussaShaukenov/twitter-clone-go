-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
