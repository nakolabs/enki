CREATE TYPE school_role AS ENUM(
    'student',
    'teacher',
    'admin',
    'head_teacher'
);

CREATE TABLE IF NOT EXISTS user_school_role (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id),
    school_id UUID REFERENCES school (id),
    role_id school_role NOT NULL DEFAULT 'student',
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
    UNIQUE (user_id, school_id, role_id, is_deleted)
)