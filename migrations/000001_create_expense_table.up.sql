CREATE TABLE IF NOT EXISTS expense(
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL,
    price FLOAT8 NOT NULL,
    created_at TIMESTAMP
);

