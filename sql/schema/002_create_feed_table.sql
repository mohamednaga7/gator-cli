-- +goose Up
CREATE TABLE feed (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP        DEFAULT now(),
    updated_at TIMESTAMP        DEFAULT now(),
    name       VARCHAR(255)  NOT NULL,
    url        VARCHAR(1000) NOT NULL UNIQUE ,
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed;