CREATE TABLE IF NOT EXISTS subject
(
    id         UUID    NOT NULL PRIMARY KEY,
    name       VARCHAR NOT NULL,
    school_id  UUID    NOT NULL REFERENCES school (id),
    created_at BIGINT  NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT  NOT NULL DEFAULT 0
)