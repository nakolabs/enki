ALTER TABLE ppdb_student 
ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS email VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'registered';

CREATE INDEX IF NOT EXISTS idx_ppdb_student_ppdb_id ON ppdb_student(ppdb_id);
CREATE INDEX IF NOT EXISTS idx_ppdb_student_email ON ppdb_student(email);
CREATE INDEX IF NOT EXISTS idx_ppdb_student_status ON ppdb_student(status);
