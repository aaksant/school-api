CREATE TABLE IF NOT EXISTS teaching_assignments (
    id INTEGER PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id),
    subject_id INTEGER NOT NULL REFERENCES subjects(id),
    teacher_id INTEGER NOT NULL REFERENCES teachers(id),
    UNIQUE (class_id, subject_id) -- 1 class, 1 subject teacher
);
