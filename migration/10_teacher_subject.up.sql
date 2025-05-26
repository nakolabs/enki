CREATE TABLE IF NOT EXISTS techer_subject
(
    id         UUID   NOT NULL PRIMARY KEY,
    techer_id   UUID   NOT NULL REFERENCES users (id),
    subject_id UUID   NOT NULL REFERENCES subject (id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL DEFAULT 0
)
