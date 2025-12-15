CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date_taken TIMESTAMP WITH TIME ZONE,
    description TEXT,
    location VARCHAR(500),
    custom_fields JSONB,
    thumbnail_photo_id INTEGER,
    photographer_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_photographer FOREIGN KEY (photographer_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_albums_photographer_id ON albums(photographer_id);
CREATE INDEX idx_albums_date_taken ON albums(date_taken);
CREATE INDEX idx_albums_custom_fields ON albums USING GIN(custom_fields);
