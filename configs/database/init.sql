CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    CHECK (char_length(username) <= 255),
    CHECK (char_length(username) > 1),
    email TEXT UNIQUE NOT NULL,
    CHECK (char_length(email) <= 255),
    CHECK (char_length(email) > 1),
    full_name TEXT,
    CHECK (char_length(full_name) <= 255),
    CHECK (char_length(full_name) > 1),
    create_date timestamp with time zone default CURRENT_TIMESTAMP,
    update_date timestamp with time zone default CURRENT_TIMESTAMP
);

CREATE TABLE concerts (
    concert_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    CHECK (char_length(title) <= 255),
    CHECK (char_length(title) > 1),
    description TEXT,
    location TEXT,
    author_id UUID NOT NULL references users(user_id),
    create_date timestamp with time zone default CURRENT_TIMESTAMP,
    update_date timestamp with time zone default CURRENT_TIMESTAMP
);

CREATE TABLE tickets (
    ticket_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    verification_token TEXT NOT NULL,
    CHECK (char_length(verification_token) <= 255),
    CHECK (char_length(verification_token) > 1),
    create_date timestamp with time zone default CURRENT_TIMESTAMP,
    update_date timestamp with time zone default CURRENT_TIMESTAMP,
    concert_id BIGINT NOT NULL references concerts(concert_id),
    user_id UUID NOT NULL references users(user_id)
)
