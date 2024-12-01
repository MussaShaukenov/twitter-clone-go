-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tweet_service_schema.tweets
(
    id         bigserial PRIMARY KEY,
    title      varchar(155)                NOT NULL,
    content    varchar(300)                NOT NULL,
    topic      varchar(64),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp                            DEFAULT current_timestamp
);
ALTER TABLE tweet_service_schema.tweets
    ADD COLUMN user_id INT NOT NULL;

ALTER TABLE tweet_service_schema.tweets
    ADD CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS tweet_service_schema.tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
CREATE INDEX tags_name_idx ON tweet_service_schema.tags (name);
CREATE TABLE tweet_service_schema.tweet_tags (
    id SERIAL PRIMARY KEY,                -- Auto-incrementing primary key
    tweet_id INT NOT NULL,                -- Foreign key to the 'tweets' table
    tag_id INT NOT NULL,                  -- Foreign key to the 'tags' table
    CONSTRAINT fk_tweet FOREIGN KEY (tweet_id) REFERENCES tweets (id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    CONSTRAINT unique_tweet_tag UNIQUE (tweet_id, tag_id)  -- Ensure unique tweet-tag pairs
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tweet_service_schema.tweets;
ALTER TABLE tweet_service_schema.tweets
    DROP CONSTRAINT IF EXISTS fk_user;

ALTER TABLE tweet_service_schema.tweets
    DROP COLUMN IF EXISTS user_id;

DROP TABLE IF EXISTS tweet_service_schema.tags;
DROP INDEX IF EXISTS tweet_service_schema.tags_name_idx;
DROP TABLE IF EXISTS tweet_service_schema.tweet_tags;
-- +goose StatementEnd
