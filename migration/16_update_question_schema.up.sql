ALTER TABLE question 
ADD COLUMN question_type VARCHAR(20) NOT NULL DEFAULT 'essay' CHECK (question_type IN ('multiple_choice', 'essay')),
ADD COLUMN options JSONB NULL,
ADD COLUMN correct_answer VARCHAR(10) NULL;

-- Update existing questions to be essay type
UPDATE question SET question_type = 'essay' WHERE question_type = 'essay';

-- Add comment for clarity
COMMENT ON COLUMN question.correct_answer IS 'Correct option ID for multiple choice questions, NULL for essay questions';
