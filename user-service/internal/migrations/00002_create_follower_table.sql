-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers
(
    id         SERIAL PRIMARY KEY,
    follower_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    followee_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers;
-- +goose StatementEnd
