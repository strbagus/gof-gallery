-- +goose Up
CREATE TABLE IF NOT EXISTS event_access_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(64) UNIQUE NOT NULL,
    event_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_event_tokens_value ON event_access_tokens(token);

-- +goose Down
DROP INDEX IF EXISTS idx_event_tokens_value;
DROP TABLE IF EXISTS event_access_tokens;
