-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_service_schema.users
(
    id         SERIAL PRIMARY KEY,                  -- Auto-incrementing primary key
    first_name VARCHAR(255)        NOT NULL,        -- First name, not null
    last_name  VARCHAR(255)        NOT NULL,        -- Last name, not null
    email      VARCHAR(255) UNIQUE NOT NULL,        -- Email, unique and not null
    age        INT CHECK (age >= 0),                -- Age with a check constraint (non-negative)
    username   VARCHAR(255) UNIQUE NOT NULL,        -- Username, unique and not null
    password   VARCHAR(255)        NOT NULL,        -- Password, hashed, not null
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Created timestamp
    updated_at timestamp DEFAULT current_timestamp
);
CREATE TABLE user_service_schema.sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE user_service_schema.users
ADD COLUMN is_first_login BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_service_schema.users;
DROP TABLE IF EXISTS user_service_schema.sessions;
ALTER TABLE user_service_schema.users DROP COLUMN is_first_login;
-- +goose StatementEnd
