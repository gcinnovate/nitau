DROP DATABASE nitau;

CREATE DATABASE nitau;

ALTER DATABASE nitau OWNER TO postgres;

\connect nitau

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    email TEXT DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO users (username, password, name, email, phone, is_active)
    VALUES('api', crypt('test1234', gen_salt('bf')), 'API User', 'sekiskylink@gmail.com', '+256782820208', TRUE);
