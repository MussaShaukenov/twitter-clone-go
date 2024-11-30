-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN is_first_login BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_first_login;
-- +goose StatementEnd
