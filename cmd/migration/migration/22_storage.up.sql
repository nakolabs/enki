CREATE TABLE storage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    folder VARCHAR(255),
    width INTEGER,
    height INTEGER,
    format VARCHAR(50),
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
);

-- Add indexes

CREATE INDEX idx_storage_public_id ON storage(public_id);

CREATE INDEX idx_storage_file_type ON storage(file_type);

CREATE INDEX idx_storage_created_at ON storage(created_at);