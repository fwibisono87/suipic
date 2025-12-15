ALTER TABLE albums
    ADD CONSTRAINT fk_thumbnail_photo FOREIGN KEY (thumbnail_photo_id) REFERENCES photos(id) ON DELETE SET NULL;

CREATE INDEX idx_albums_thumbnail_photo_id ON albums(thumbnail_photo_id);
