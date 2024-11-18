-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tweets
(
    id         bigserial PRIMARY KEY,
    title      varchar(155)                NOT NULL,
    content    varchar(300)                NOT NULL,
    topic      varchar(64),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp                            DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tweets;
-- +goose StatementEnd