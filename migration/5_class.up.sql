CREATE TABLE IF NOT EXISTS class
(
    id         UUID    NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id  UUID    NOT NULL REFERENCES school (id),
    name       VARCHAR NOT NULL,
    created_at BIGINT  NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT  NOT NULL DEFAULT 0
)