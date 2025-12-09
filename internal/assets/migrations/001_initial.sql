-- +migrate Up
CREATE TABLE short_urls(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY
    long_url string NOT NULL
    short_url string UNIQUE
)

CREATE INDEX short_url_idx on short_urls(short_url)
-- +migrate Down
