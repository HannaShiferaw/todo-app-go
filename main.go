// Command todo-app-go is a todo list: a Go + Gin API backed by a
// SQLite file, with a small web page in public/ that talks to it.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"todo-app-go/database"
	"todo-app-go/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// corsMiddleware allows the web page to call the API from any origin
// during local development, and answers OPTIONS preflight requests.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// setupRouter wires every route. It is separate from main so tests
// can build a router against a temporary database.
func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	r.Static("/public", "./public")
	r.StaticFile("/", "./public/index.html")

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is running"})
	})
	r.POST("/api/todos", handlers.CreateTodo(db))
	r.GET("/api/todos", handlers.ListTodos(db))
	r.GET("/api/todos/:id", handlers.GetTodo(db))
	r.PUT("/api/todos/:id", handlers.UpdateTodo(db))
	r.DELETE("/api/todos/:id", handlers.DeleteTodo(db))
	return r
}

// port returns the PORT environment variable, or 8080 when unset.
func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8080"
}

func main() {
	database.Connect()
	router := setupRouter(database.DB)

	server := &http.Server{
		Addr:         ":" + port(),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("Starting server at http://localhost:%s", port())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
