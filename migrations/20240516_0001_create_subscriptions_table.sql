CREATE TABLE IF NOT EXISTS frequencies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO frequencies (name) VALUES ('Hourly'), ('Daily')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS subscriptions (
    email TEXT PRIMARY KEY,
    city TEXT NOT NULL,
    frequency INTEGER NOT NULL REFERENCES frequencies(id),
    token TEXT NOT NULL UNIQUE,
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    subscribed_at TIMESTAMPTZ DEFAULT now(),
    last_sent_at TIMESTAMPTZ DEFAULT
);
