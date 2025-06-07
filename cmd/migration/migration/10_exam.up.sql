CREATE TABLE IF NOT EXISTS exam (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    school_id UUID NOT NULL REFERENCES school (id),
    subject_id UUID NOT NULL REFERENCES subject (id),
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_at BIGINT NOT NULL DEFAULT (
        EXTRACT(
            EPOCH
            FROM
                now()
        ) * 1000
    ) :: BIGINT,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_at BIGINT NOT NULL DEFAULT 0,
    updated_by UUID NOT NULL REFERENCES users(id),
    deleted_at BIGINT DEFAULT 0,
    deleted_by UUID DEFAULT NULL REFERENCES users(id),
    UNIQUE (name, school_id, subject_id, is_deleted)
)