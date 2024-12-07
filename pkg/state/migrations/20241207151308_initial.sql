CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    email text NOT NULL,
    password_hash text NOT NULL
)