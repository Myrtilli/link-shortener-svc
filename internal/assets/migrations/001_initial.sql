-- +migrate Up
CREATE TABLE short_urls(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url TEXT UNIQUE NOT NULL
);

CREATE INDEX short_url_idx on short_urls(short_url);
-- +migrate Down
DROP TABLE short_urls;