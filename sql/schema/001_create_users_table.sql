-- +goose Up
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP        DEFAULT now(),
    updated_at TIMESTAMP        DEFAULT now(),
    name       VARCHAR(255) NOT NULL
);

CREATE TABLE feed (
                      id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      created_at TIMESTAMP        DEFAULT now(),
                      updated_at TIMESTAMP        DEFAULT now(),
                      name       VARCHAR(255)  NOT NULL,
                      url        VARCHAR(1000) NOT NULL UNIQUE ,
                      user_id    UUID REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows (
                              id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
                              feed_id    UUID REFERENCES feed  (id) ON DELETE CASCADE,
                              created_at TIMESTAMP DEFAULT now(),
                              updated_at TIMESTAMP DEFAULT now(),
                              UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;

DROP TABLE feed;

DROP TABLE users;