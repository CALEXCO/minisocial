-- +goose Up
CREATE TABLE users (
                       id         BIGSERIAL PRIMARY KEY,
                       name   TEXT NOT NULL UNIQUE,
                       email      TEXT NOT NULL UNIQUE,
                        password TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (name, email, password) VALUES ('Alex', 'Alex@gamil.com', 'Alex123');
INSERT INTO users (name, email, password) VALUES ('Alex1', 'Alex1@gamil.com', 'Alex123');
INSERT INTO users (name, email, password) VALUES ('Alex2', 'Alex2@gamil.com', 'Alex123');
INSERT INTO users (name, email, password) VALUES ('Alex3', 'Alex@3gamil.com', 'Alex123');

-- +goose Down
DROP TABLE users;