-- For auth
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    person_id INTEGER UNIQUE,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher', 'staff', 'executive')),
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_user_updated_at
AFTER UPDATE ON users
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE users SET updated_at = datetime('now') WHERE id = NEW.id;
END;
