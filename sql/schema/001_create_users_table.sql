-- +goose Up
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP        DEFAULT now(),
    updated_at TIMESTAMP        DEFAULT now(),
    name       VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;