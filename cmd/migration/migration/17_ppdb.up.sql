CREATE TABLE IF NOT EXISTS ppdb (
    id UUID NOT NULL PRIMARY KEY,
    school_id UUID NOT NULL REFERENCES school (id),
    start_at BIGINT NOT NULL,
    end_at BIGINT NOT NULL,
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
    deleted_by UUID DEFAULT NULL REFERENCES users(id)
)