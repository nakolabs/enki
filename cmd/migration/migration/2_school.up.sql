CREATE TYPE school_level AS ENUM (
    'preschool',
    'kindergarten',
    'elementary',
    'junior',
    'senior',
    'college'
);

CREATE TABLE IF NOT EXISTS school (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR NOT NULL,
    level school_level NOT NULL,
    address VARCHAR(255) NOT NULL DEFAULT '',
    city VARCHAR(100) NOT NULL DEFAULT '',
    province VARCHAR(100) NOT NULL DEFAULT '',
    postal_code VARCHAR(20) NOT NULL DEFAULT '',
    phone VARCHAR(20) NOT NULL DEFAULT '',
    email VARCHAR(100) NOT NULL DEFAULT '',
    website VARCHAR(100) NOT NULL DEFAULT '',
    logo VARCHAR(255) NOT NULL DEFAULT '',
    banner VARCHAR(255) NOT NULL DEFAULT '',
    established_year INT NOT NULL DEFAULT 0,
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
    deleted_by UUID REFERENCES users(id)
);