ALTER TABLE question 
DROP COLUMN IF EXISTS question_type,
DROP COLUMN IF EXISTS options,
DROP COLUMN IF EXISTS correct_answer;
