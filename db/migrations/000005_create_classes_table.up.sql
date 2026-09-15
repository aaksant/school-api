CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    homeroom_teacher_id INTEGER UNIQUE REFERENCES teachers(id) ON DELETE SET NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_class_updated_at
AFTER UPDATE ON classes
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE classes SET updated_at = datetime('now') WHERE id = NEW.id;
END;
