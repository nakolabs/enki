CREATE TABLE IF NOT EXISTS exam_question
(
    id          UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    exam_id     UUID             NOT NULL REFERENCES exam (id),
    question_id UUID             NOT NULL REFERENCES question (id),
    created_at  BIGINT           NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at  BIGINT           NOT NULL DEFAULT 0
)