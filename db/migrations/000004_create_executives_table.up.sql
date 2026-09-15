CREATE TABLE IF NOT EXISTS executives (
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER IF NOT EXISTS trg_executive_updated_at
AFTER UPDATE ON executives
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE executives SET updated_at = datetime('now') WHERE id = NEW.id;
END;
