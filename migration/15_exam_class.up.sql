CREATE TABLE IF NOT EXISTS exam_class
(
    id         UUID   NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id    UUID   NOT NULL REFERENCES exam (id),
    class_id   UUID   NOT NULL REFERENCES class (id),
    created_at BIGINT NOT NULL             DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL             DEFAULT 0,
    UNIQUE(exam_id, class_id)
);
