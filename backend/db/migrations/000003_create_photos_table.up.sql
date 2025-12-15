CREATE TYPE pick_reject_state AS ENUM ('none', 'pick', 'reject');

CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    album_id INTEGER NOT NULL,
    filename VARCHAR(500) NOT NULL,
    title VARCHAR(255),
    date_time TIMESTAMP WITH TIME ZONE,
    exif_data JSONB,
    pick_reject_state pick_reject_state NOT NULL DEFAULT 'none',
    stars INTEGER DEFAULT 0 CHECK (stars >= 0 AND stars <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_album FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
);

CREATE INDEX idx_photos_album_id ON photos(album_id);
CREATE INDEX idx_photos_date_time ON photos(date_time);
CREATE INDEX idx_photos_pick_reject_state ON photos(pick_reject_state);
CREATE INDEX idx_photos_stars ON photos(stars);
CREATE INDEX idx_photos_exif_data ON photos USING GIN(exif_data);
