CREATE TABLE IF NOT EXISTS teacher_subject (
    id UUID NOT NULL PRIMARY KEY,
    teacher_id UUID NOT NULL REFERENCES users (id),
    subject_id UUID NOT NULL REFERENCES subject (id),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at BIGINT NOT NULL DEFAULT (
        EXTRACT(
            EPOCH
            FROM
                now()
        ) * 1000
    ) :: BIGINT,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_at BIGINT NOT NULL DEFAULT 0,
    updated_by UUID REFERENCES users(id),
    deleted_at BIGINT DEFAULT 0,
    deleted_by UUID DEFAULT NULL REFERENCES users(id),
    UNIQUE (teacher_id, subject_id, is_deleted)
)