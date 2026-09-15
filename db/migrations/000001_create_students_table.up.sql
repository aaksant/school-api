CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    date_of_birth TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_student_updated_at
AFTER UPDATE ON students
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE students SET updated_at = datetime('now') WHERE id = NEW.id;
END;
