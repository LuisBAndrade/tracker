-- +goose Up
CREATE TABLE refresh_tokens (
    token VARCHAR(255) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;