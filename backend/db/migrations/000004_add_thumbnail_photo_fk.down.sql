DROP INDEX IF EXISTS idx_albums_thumbnail_photo_id;
ALTER TABLE albums DROP CONSTRAINT IF EXISTS fk_thumbnail_photo;
