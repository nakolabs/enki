DROP INDEX IF EXISTS idx_question_school_subject;
DROP INDEX IF EXISTS idx_question_type;

ALTER TABLE question 
DROP COLUMN IF EXISTS school_id,
DROP COLUMN IF EXISTS subject_id,
DROP COLUMN IF EXISTS difficulty_level,
DROP COLUMN IF EXISTS points;
