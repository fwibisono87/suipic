CREATE TABLE album_users (
    id SERIAL PRIMARY KEY,
    album_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_album_users_album FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
    CONSTRAINT fk_album_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_album_user UNIQUE (album_id, user_id)
);

CREATE INDEX idx_album_users_album_id ON album_users(album_id);
CREATE INDEX idx_album_users_user_id ON album_users(user_id);
