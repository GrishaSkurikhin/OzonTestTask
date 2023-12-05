CREATE TABLE urls (
	long_url VARCHAR NOT NULL,
    short_url VARCHAR NOT NULL
);

CREATE INDEX urls_short_url_idx ON urls (short_url);

ALTER TABLE urls ADD CONSTRAINT unique_urls UNIQUE (long_url, short_url);
