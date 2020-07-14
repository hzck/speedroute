CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL,
    last_updated TIMESTAMP NOT NULL
);