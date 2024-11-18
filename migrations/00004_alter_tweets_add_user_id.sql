-- +goose Up
-- +goose StatementBegin
ALTER TABLE tweets
    ADD COLUMN user_id INT NOT NULL;

ALTER TABLE tweets
    ADD CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tweets
    DROP CONSTRAINT IF EXISTS fk_user;

ALTER TABLE tweets
    DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd
