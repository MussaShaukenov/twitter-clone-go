-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    first_name varchar(32) NOT NULL,
    last_name varchar(32) NOT NULL,
    email varchar(50) NOT NULL,
    age int,
    password_hash varchar NOT NULL,
    activated bool DEFAULT false,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
