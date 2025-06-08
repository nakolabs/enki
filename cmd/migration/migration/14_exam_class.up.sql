CREATE TABLE IF NOT EXISTS exam_class (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    exam_id UUID NOT NULL REFERENCES exam (id),
    class_id UUID NOT NULL REFERENCES class (id),
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
    UNIQUE(exam_id, class_id, is_deleted)
);