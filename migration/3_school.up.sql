CREATE TYPE school_level AS ENUM ('preschool', 'kindergarten', 'elementary', 'junior', 'senior', 'college');
CREATE TABLE IF NOT EXISTS school
(
    id         UUID         NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR      NOT NULL,
    level      school_level NOT NULL,
    created_at BIGINT       NOT NULL             DEFAULT (EXTRACT(EPOCH FROM now()) * 1000)::BIGINT,
    updated_at BIGINT       NOT NULL             DEFAULT 0
);