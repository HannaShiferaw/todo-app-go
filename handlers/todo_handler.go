// Package handlers holds the HTTP handlers for the todo API.
// Each function takes the database and returns a Gin handler, so
// the same code serves real requests and tests.
package handlers

import (
	"net/http"
	"strings"

	"todo-app-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var validPriorities = map[string]bool{"low": true, "medium": true, "high": true}
var validCategories = map[string]bool{
	"work": true, "personal": true, "shopping": true, "health": true, "other": true,
}

// normalize fills in defaults and checks priority and category.
// It returns false (after writing the error) when the todo is invalid.
func normalize(c *gin.Context, todo *models.Todo) bool {
	todo.Title = strings.TrimSpace(todo.Title)
	if todo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return false
	}
	if todo.Priority == "" {
		todo.Priority = "medium"
	} else if !validPriorities[todo.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priority must be low, medium or high"})
		return false
	}
	if todo.Category == "" {
		todo.Category = "personal"
	} else if !validCategories[todo.Category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category must be work, personal, shopping, health or other"})
		return false
	}
	return true
}

// findTodo loads a todo by its :id, or writes 404 and returns false.
func findTodo(db *gorm.DB, c *gin.Context) (models.Todo, bool) {
	var todo models.Todo
	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return todo, false
	}
	return todo, true
}

// CreateTodo handles POST /api/todos.
func CreateTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !normalize(c, &todo) {
			return
		}
		if err := db.Create(&todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, todo)
	}
}

// ListTodos handles GET /api/todos.
func ListTodos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var todos []models.Todo
		if err := db.Find(&todos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if todos == nil {
			todos = []models.Todo{}
		}
		c.JSON(http.StatusOK, todos)
	}
}

// GetTodo handles GET /api/todos/:id.
func GetTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo, ok := findTodo(db, c)
		if !ok {
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

// UpdateTodo handles PUT /api/todos/:id. It replaces the fields
// sent in the request; a missing title is rejected.
func UpdateTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo, ok := findTodo(db, c)
		if !ok {
			return
		}
		var updated models.Todo
		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !normalize(c, &updated) {
			return
		}
		todo.Title = updated.Title
		todo.Description = updated.Description
		todo.Completed = updated.Completed
		todo.Priority = updated.Priority
		todo.DueDate = updated.DueDate
		todo.Category = updated.Category
		if err := db.Save(&todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

// DeleteTodo handles DELETE /api/todos/:id.
func DeleteTodo(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		todo, ok := findTodo(db, c)
		if !ok {
			return
		}
		if err := db.Delete(&todo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
