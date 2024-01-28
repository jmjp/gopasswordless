-- Create schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS public;

-- Create the User table in the PostgreSQL schema
CREATE TABLE IF NOT EXISTS Users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username STRING NOT NULL UNIQUE,
    email STRING NOT NULL UNIQUE,
    avatar STRING,
    blocked BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT current_timestamp(),
    updated_at TIMESTAMPTZ DEFAULT current_timestamp()
);

-- Create the MagicLink table in the PostgreSQL schema
CREATE TABLE IF NOT EXISTS Magic_Links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code STRING NOT NULL,
    user_id UUID NOT NULL REFERENCES Users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    cookie STRING NOT NULL UNIQUE,
    valid_until TIMESTAMPTZ NOT NULL,
	used BOOLEAN NOT NULL DEFAULT false
);


-- Create the Session table in the PostgresSQL Schema
CREATE TABLE IF NOT EXISTS Sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES Users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    valid_until TIMESTAMPTZ NOT NULL,
    user_Agent STRING,
    ip STRING,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    city STRING,
    region STRING,
    country STRING,
    isp STRING,
    created_at TIMESTAMPTZ DEFAULT current_timestamp(),
    updated_at TIMESTAMPTZ DEFAULT current_timestamp()
);


CREATE INDEX IF NOT EXISTS magic_links_code_storing_rec_idx ON magic_links (code) STORING (user_id, cookie, valid_until, used); 
CREATE INDEX IF NOT EXISTS magic_links_code_cookie_token_valid_until_storing_rec_idx ON magic_links (code, cookie, valid_until) STORING (user_id, used);