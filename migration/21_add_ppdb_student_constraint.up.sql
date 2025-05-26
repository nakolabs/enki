-- Add unique constraint to prevent duplicate registration
ALTER TABLE ppdb_student ADD CONSTRAINT unique_ppdb_student UNIQUE (ppdb_id, student_id);

-- Add index for better performance
CREATE INDEX IF NOT EXISTS idx_ppdb_student_user_id ON ppdb_student(student_id);
