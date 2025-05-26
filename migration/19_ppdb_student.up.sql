CREATE TABLE IF NOT EXISTS ppdb_student
(
    id         UUID   NOT NULL PRIMARY KEY,
    ppdb_id    UUID   NOT NULL REFERENCES ppdb (id),
    student_id UUID   NOT NULL REFERENCES users (id),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
)