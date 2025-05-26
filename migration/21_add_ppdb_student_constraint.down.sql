DROP INDEX IF EXISTS idx_ppdb_student_user_id;
ALTER TABLE ppdb_student DROP CONSTRAINT IF EXISTS unique_ppdb_student;
