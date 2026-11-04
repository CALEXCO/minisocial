-- +goose Up
CREATE TABLE users (
                       id         BIGSERIAL PRIMARY KEY,
                       name   TEXT NOT NULL UNIQUE,
                       email      TEXT NOT NULL UNIQUE,
                        password TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;