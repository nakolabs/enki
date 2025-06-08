CREATE TABLE IF NOT EXISTS class_student (
    id UUID NOT NULL PRIMARY KEY,
    student_id UUID NOT NULL REFERENCES users (id),
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
    UNIQUE (student_id, class_id, is_deleted)
)