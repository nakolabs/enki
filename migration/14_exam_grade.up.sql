CREATE TABLE IF NOT EXISTS exam_grade
(
    id         UUID   NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id    UUID   NOT NULL REFERENCES exam (id),
    student_id UUID   NOT NULL REFERENCES users (id),
    grade      DECIMAL(5,2),
    answers    TEXT,
    created_at BIGINT NOT NULL             DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL             DEFAULT 0,
    UNIQUE(exam_id, student_id)
);