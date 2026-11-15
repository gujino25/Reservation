CREATE TABLE IF NOT EXISTS guests (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    phone TEXT UNIQUE,
    created_at TIMESTAMPTZ
);
