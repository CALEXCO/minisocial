-- +goose Up
CREATE TABLE users (
                       id         BIGSERIAL PRIMARY KEY,
                       username   TEXT NOT NULL UNIQUE,
                       email      TEXT NOT NULL UNIQUE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;