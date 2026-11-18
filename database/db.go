// Package database opens the SQLite database and creates the
// todos table if it does not exist yet. No separate setup step is
// needed: the table is created automatically on startup.
//
// The database file is TODO_DB_PATH, or todo.db by default.
// There is no server, username or password to configure —
// SQLite is just a file on disk.
package database

import (
	"log"
	"os"
	"strings"

	"todo-app-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the shared database connection, set by Connect.
var DB *gorm.DB

// dbPath returns the database file path from the environment.
func dbPath() string {
	if p := strings.TrimSpace(os.Getenv("TODO_DB_PATH")); p != "" {
		return p
	}
	return "todo.db"
}

// ConnectWithPath opens the database at path and migrates the schema.
// It is separate from Connect so tests can use a temporary database.
func ConnectWithPath(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		return nil, err
	}
	return db, nil
}

// Connect opens the default database and stores it in DB.
// It stops the program if the database cannot be opened.
func Connect() {
	db, err := ConnectWithPath(dbPath())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
	log.Println("Database connected successfully")
}
