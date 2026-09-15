CREATE TABLE IF NOT EXISTS subjects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_subject_updated_at
AFTER UPDATE ON subjects
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE subjects SET updated_at = datetime('now') WHERE id = NEW.id;
END;
