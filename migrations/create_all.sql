-- +goose Up
CREATE TABLE posts (
                       id         BIGSERIAL PRIMARY KEY,
                       user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       content    TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);

-- +goose Down
DROP TABLE posts;