package main

import (
	"net/http"
	"todo-app-go/database"
	"todo-app-go/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()


	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})


	r.Static("/public", "./public")
	r.StaticFile("/", "./public/index.html")

	
	database.Connect()

	
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})

	// Create a new todo
	r.POST("/api/todos", func(c *gin.Context) {
		var todo models.Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := database.DB.Create(&todo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	// Get all todos
	r.GET("/api/todos", func(c *gin.Context) {
		var todos []models.Todo
		result := database.DB.Find(&todos)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		if todos == nil {
			todos = []models.Todo{}
		}
		c.JSON(http.StatusOK, todos)
	})

	// Get a specific todo
	r.GET("/api/todos/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		result := database.DB.First(&todo, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Todo not found",
			})
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	// Update a todo
	r.PUT("/api/todos/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		if err := database.DB.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		var updatedTodo models.Todo
		if err := c.ShouldBindJSON(&updatedTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		todo.Title = updatedTodo.Title
		todo.Description = updatedTodo.Description
		todo.Completed = updatedTodo.Completed
		todo.Priority = updatedTodo.Priority
		todo.DueDate = updatedTodo.DueDate
		todo.Category = updatedTodo.Category
		database.DB.Save(&todo)
		c.JSON(http.StatusOK, todo)
	})

	// Delete a todo
	r.DELETE("/api/todos/:id", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		if err := database.DB.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Todo not found",
			})
			return
		}
		database.DB.Delete(&todo)
		c.JSON(http.StatusOK, gin.H{
			"message": "Todo deleted successfully",
		})
	})

	r.Run()
}
