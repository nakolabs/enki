CREATE TABLE IF NOT EXISTS user_school_role
(
    id         UUID   NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID REFERENCES users (id),
    school_id  UUID REFERENCES school (id),
    role_id    UUID REFERENCES role (id),
    created_at BIGINT NOT NULL             DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT NOT NULL             DEFAULT 0,
    UNIQUE (user_id, school_id, role_id)
)