CREATE TABLE IF NOT EXISTS adv (
    ad_id text PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    photos text[] NOT NULL,
    cost numeric NOT NULL
);