DROP INDEX IF EXISTS idx_ppdb_student_status;
DROP INDEX IF EXISTS idx_ppdb_student_email;
DROP INDEX IF EXISTS idx_ppdb_student_ppdb_id;

ALTER TABLE ppdb_student 
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS name;
