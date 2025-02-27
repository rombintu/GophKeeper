BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hex_keys BYTEA,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE secrets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    secret_type INTEGER NOT NULL CHECK (secret_type BETWEEN 0 AND 3) DEFAULT 0, -- Ограничение 0-3
    user_email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    version NUMERIC(10,2) DEFAULT 0,
    payload BYTEA,
    FOREIGN KEY (user_email) REFERENCES users(email)
);

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    last_check TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    alive BOOLEAN DEFAULT true
);

COMMIT;