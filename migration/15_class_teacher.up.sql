CREATE TABLE IF NOT EXISTS class_teacher
(
    id         UUID   NOT NULL PRIMARY KEY,
    teacher_id UUID   NOT NULL REFERENCES users (id),
    class_id   UUID   NOT NULL REFERENCES class (id),
    created_at BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL DEFAULT 0
)