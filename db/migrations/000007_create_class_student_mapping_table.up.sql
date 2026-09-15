CREATE TABLE IF NOT EXISTS class_student_mapping (
    id INTEGER PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id),
    student_id INTEGER NOT NULL REFERENCES students(id),
    UNIQUE(student_id)
);
