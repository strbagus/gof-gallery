-- +goose Up
CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL,
    event_slug VARCHAR(100) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    preview VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_event
        FOREIGN KEY(event_id) 
        REFERENCES events(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_event_photo 
        UNIQUE(event_id, filename)
);
CREATE INDEX idx_photos_event_id ON photos(event_id);

-- +goose Down
DROP TABLE IF EXISTS photos;
