CREATE TABLE if not exists urls
(
    id           SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url    TEXT NOT NULL UNIQUE,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_short_url ON urls(short_url);
