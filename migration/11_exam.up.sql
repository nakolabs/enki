CREATE TABLE IF NOT EXISTS exam
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name       VARCHAR          NOT NULL,
    school_id  UUID             NOT NULL REFERENCES school (id),
    subject_id UUID             NOT NULL REFERENCES subject (id),
    created_at BIGINT           NOT NULL DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT           NOT NULL DEFAULT 0
)