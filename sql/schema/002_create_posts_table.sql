-- +goose Up
CREATE TABLE posts
(
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at   TIMESTAMP        DEFAULT now(),
    updated_at   TIMESTAMP        DEFAULT now(),
    published_at TIMESTAMP        DEFAULT now(),
    title        VARCHAR(255)  NOT NULL,
    url          VARCHAR(1000) NOT NULL UNIQUE,
    description  TEXT,
    feed_id      UUID REFERENCES feed (id)
);

-- +goose Down
DROP TABLE posts;