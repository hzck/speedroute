CREATE TABLE websession (
    token UUID PRIMARY KEY,
    account_id INTEGER REFERENCES account (id) ON DELETE CASCADE NOT NULL,
    expire_at TIMESTAMPTZ NOT NULL
);