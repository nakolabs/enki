CREATE TABLE IF NOT EXISTS ppdb
(
    id         UUID   NOT NULL PRIMARY KEY,
    school_id  UUID   NOT NULL REFERENCES school (id),
    start_at   BIGINT NOT NULL,
    end_at     BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
)