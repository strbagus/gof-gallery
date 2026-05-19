-- +goose Up
-- 1. Create the Events Table
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    location VARCHAR(255),
    date DATE,
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    salt VARCHAR(32) NOT NULL DEFAULT md5(random()::text),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Add Performance Index
CREATE INDEX idx_events_slug ON events(slug);

-- +goose Down
DROP TABLE IF EXISTS events;
