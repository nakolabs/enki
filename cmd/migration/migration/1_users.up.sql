CREATE TYPE gender AS ENUM ('male', 'female', '');
CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    role user_role NOT NULL DEFAULT 'user',
    phone VARCHAR(20) NOT NULL DEFAULT '',
    date_of_birth DATE NOT NULL DEFAULT '1970-01-01',
    gender gender NOT NULL DEFAULT '',
    address VARCHAR(255) NOT NULL DEFAULT '',
    city VARCHAR(100) NOT NULL DEFAULT '',
    country VARCHAR(100) NOT NULL DEFAULT '',
    avatar VARCHAR(255) NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    parent_name VARCHAR(100) NOT NULL DEFAULT '',
    parent_phone VARCHAR(20) NOT NULL DEFAULT '',
    parent_email VARCHAR(100) NOT NULL DEFAULT '',
    is_deleted BOOLEAN NOT NULL DEFAULT false,
    created_at BIGINT NOT NULL DEFAULT (
        EXTRACT(
            EPOCH
            FROM
                now()
        ) * 1000
    ) :: BIGINT,
    updated_at BIGINT NOT NULL DEFAULT 0,
    deleted_at BIGINT DEFAULT 0,
    deleted_by UUID DEFAULT NULL REFERENCES users(id),
    UNIQUE (email, is_deleted)
);