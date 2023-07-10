CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    displayname VARCHAR(30) NOT NULL,
    password VARCHAR NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    last_updated TIMESTAMPTZ NOT NULL
);