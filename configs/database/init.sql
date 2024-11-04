CREATE TABLE users (
    user_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username TEXT UNIQUE NOT NULL,
    CHECK (char_length(uesrname) <= 255),
    CHECK (char_length(username) > 1),
    email TEXT UNIQUE NOT NULL,
    CHECK (char_length(email) <= 255),
    CHECK (char_length(email) > 1),
    full_name TEXT,
    CHECK (char_length(full_name) <= 255),
    CHECK (char_length(full_name) > 1),
    create_date timestamp with time zone NOT NULL,
    update_date timestamp with time zone
);

CREATE TABLE concerts (
    concert_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    CHECK (char_length(title) <= 255),
    CHECK (char_length(title) > 1),
    description TEXT,
    location TEXT,
    create_date timestamp with time zone NOT NULL,
    update_date timestamp with time zone
);

CREATE TABLE tickets (
    ticket_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    verification_token TEXT NOT NULL,
    CHECK (char_length(verification_token) <= 255),
    CHECK (char_length(verification_token) > 1),
    create_date timestamp with time zone NOT NULL,
    update_date timestamp with time zone,
    concert_id BIGINT refernces concerts(concert_id),
    user_id BIGINT references users(user_id)
)