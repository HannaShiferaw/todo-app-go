package database

import (
	"log"
	"todo-app-go/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


var DB *gorm.DB


func Connect() {
	
	dsn := "todo.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = db
	log.Println("Database connected successfully")
	DB.AutoMigrate(&models.Todo{})
	log.Println("Database migrated successfully")
}
