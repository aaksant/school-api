CREATE TABLE IF NOT EXISTS staffs (
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    position TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_staff_updated_at
AFTER UPDATE ON staffs
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE staffs SET updated_at = datetime('now') WHERE id = NEW.id;
END;
