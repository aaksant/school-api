-- For auth
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL,
    student_id INTEGER UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    teacher_id INTEGER UNIQUE REFERENCES teachers(id) ON DELETE CASCADE,
    staff_id INTEGER UNIQUE REFERENCES staffs(id) ON DELETE CASCADE,
    executive_id INTEGER UNIQUE REFERENCES executives(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher', 'staff', 'executive')),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    CHECK (
        (role = 'student') = (student_id IS NOT NULL) AND
        (role = 'teacher') = (teacher_id IS NOT NULL) AND
        (role = 'staff') = (staff_id IS NOT NULL) AND
        (role = 'executive') = (executive_id IS NOT NULL)
    )
);

CREATE TRIGGER IF NOT EXISTS trg_user_updated_at
AFTER UPDATE ON users
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE users SET updated_at = datetime('now') WHERE id = NEW.id;
END;
