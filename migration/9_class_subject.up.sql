CREATE TABLE IF NOT EXISTS class_subject
(
    id         UUID   NOT NULL PRIMARY KEY,
    class_id   UUID   NOT NULL REFERENCES class (id),
    subject_id UUID   NOT NULL REFERENCES subject (id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL DEFAULT 0
)
