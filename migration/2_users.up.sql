CREATE TABLE IF NOT EXISTS users
(
    id          UUID    NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email       VARCHAR NOT NULL UNIQUE,
    name        VARCHAR NOT NULL,
    password    VARCHAR NOT NULL,
    is_verified BOOLEAN NOT NULL             DEFAULT false,
    created_at  BIGINT  NOT NULL             DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at  BIGINT  NOT NULL             DEFAULT 0
);