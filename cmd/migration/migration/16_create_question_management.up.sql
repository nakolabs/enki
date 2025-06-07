-- Add school_id and subject_id to question table for better organization
ALTER TABLE question 
ADD COLUMN school_id UUID REFERENCES school(id),
ADD COLUMN subject_id UUID REFERENCES subject(id),
ADD COLUMN difficulty_level VARCHAR(20) DEFAULT 'medium' CHECK (difficulty_level IN ('easy', 'medium', 'hard')),
ADD COLUMN points INTEGER DEFAULT 1 CHECK (points > 0);

-- Create index for better query performance
CREATE INDEX idx_question_school_subject ON question(school_id, subject_id);
CREATE INDEX idx_question_type ON question(question_type);
