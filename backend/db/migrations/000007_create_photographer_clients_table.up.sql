CREATE TABLE IF NOT EXISTS photographer_clients (
    id SERIAL PRIMARY KEY,
    photographer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(photographer_id, client_id)
);

CREATE INDEX IF NOT EXISTS idx_photographer_clients_photographer_id ON photographer_clients(photographer_id);
CREATE INDEX IF NOT EXISTS idx_photographer_clients_client_id ON photographer_clients(client_id);
