ALTER TABLE
    teacher_subject
ADD
    CONSTRAINT u_teacher_subject UNIQUE (teacher_id, subject_id);