CREATE TABLE storage_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
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
    updated_at BIGINT NOT NULL DEFAULT 0
);

-- Add indexes
CREATE INDEX idx_storage_log_user_id ON storage_log(user_id);
CREATE INDEX idx_storage_log_public_id ON storage_log(public_id);
CREATE INDEX idx_storage_log_file_type ON storage_log(file_type);
CREATE INDEX idx_storage_log_created_at ON storage_log(created_at);

-- Add foreign key constraint
ALTER TABLE storage_log 
ADD CONSTRAINT fk_storage_log_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
