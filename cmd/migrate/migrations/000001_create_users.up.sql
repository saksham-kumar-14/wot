CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username varchar(255) NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    about text NOT NULL DEFAULT '',
    friends bigint[] NOT NULL DEFAULT '{}',
    friends_of bigint[] NOT NULL DEFAULT '{}',
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
