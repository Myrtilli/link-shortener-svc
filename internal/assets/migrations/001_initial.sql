-- +migrate Up
CREATE TABLE short_urls(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url TEXT UNIQUE
);

CREATE INDEX short_url_idx on short_urls(short_url);
-- +migrate Down
DROP TABLE short_urls;