CREATE TABLE IF NOT EXISTS advertisement (
    id text PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    photos text[] NOT NULL,
    cost numeric NOT NULL
);