-- Schema for todo-app-go (SQLite).
-- You do not need to run this by hand: the app creates the table
-- automatically on startup with GORM AutoMigrate. It is kept here
-- so the shape of the data is visible without reading Go code.
CREATE TABLE IF NOT EXISTS todos (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    completed   BOOLEAN NOT NULL DEFAULT 0,
    priority    TEXT NOT NULL DEFAULT 'medium',  -- low, medium, high
    due_date    DATETIME,
    category    TEXT NOT NULL DEFAULT 'personal' -- work, personal, shopping, health, other
);
CREATE INDEX IF NOT EXISTS idx_todos_deleted_at ON todos (deleted_at);
