CREATE TABLE if not exists urls
(
    id           SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url    TEXT NOT NULL UNIQUE
);

ALTER TABLE urls
    ADD COLUMN create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

