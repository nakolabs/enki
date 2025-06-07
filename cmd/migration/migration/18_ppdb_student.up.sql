CREATE TABLE IF NOT EXISTS ppdb_student (
    id UUID NOT NULL PRIMARY KEY,
    ppdb_id UUID NOT NULL REFERENCES ppdb (id),
    student_id UUID NOT NULL REFERENCES users (id),
    status VARCHAR(50) NOT NULL DEFAULT 'registered',
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
    updated_by UUID NOT NULL REFERENCES users(id),
    deleted_at BIGINT DEFAULT 0,
    deleted_by UUID DEFAULT NULL REFERENCES users(id),
    UNIQUE (ppdb_id, student_id, is_deleted)
)