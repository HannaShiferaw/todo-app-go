package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todo-app-go/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func itoa(n int) string { return strconv.Itoa(n) }

func testRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	// A fresh in-memory database per test, so tests do not see each other's rows.
	db, err := gorm.Open(sqlite.Open("file:memdb_"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Todo{}); err != nil {
		t.Fatalf("migrating test database: %v", err)
	}
	r := gin.New()
	r.POST("/api/todos", CreateTodo(db))
	r.GET("/api/todos", ListTodos(db))
	r.GET("/api/todos/:id", GetTodo(db))
	r.PUT("/api/todos/:id", UpdateTodo(db))
	r.DELETE("/api/todos/:id", DeleteTodo(db))
	return r, db
}

func doRequest(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encoding body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestCreateAndGetTodo(t *testing.T) {
	r, _ := testRouter(t)
	rec := doRequest(t, r, http.MethodPost, "/api/todos", map[string]any{"title": "Buy milk", "priority": "high"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	var created models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decoding created todo: %v", err)
	}
	if created.ID == 0 || created.Title != "Buy milk" {
		t.Errorf("created todo = %+v", created)
	}

	rec = doRequest(t, r, http.MethodGet, "/api/todos", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d", rec.Code)
	}
	var list []models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decoding list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("list has %d todos, want 1", len(list))
	}
}

func TestCreateTodoValidation(t *testing.T) {
	r, _ := testRouter(t)

	rec := doRequest(t, r, http.MethodPost, "/api/todos", map[string]any{"title": "  "})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("empty title status = %d, want 400", rec.Code)
	}

	rec = doRequest(t, r, http.MethodPost, "/api/todos", map[string]any{"title": "x", "priority": "urgent"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("bad priority status = %d, want 400", rec.Code)
	}
}

func TestGetMissingTodoIs404(t *testing.T) {
	r, _ := testRouter(t)
	rec := doRequest(t, r, http.MethodGet, "/api/todos/9999", nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestUpdateAndDeleteTodo(t *testing.T) {
	r, _ := testRouter(t)
	rec := doRequest(t, r, http.MethodPost, "/api/todos", map[string]any{"title": "Old title"})
	var created models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decoding created todo: %v", err)
	}
	path := "/api/todos/" + itoa(int(created.ID))

	rec = doRequest(t, r, http.MethodPut, path, map[string]any{"title": "New title", "completed": true})
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d: %s", rec.Code, rec.Body.String())
	}
	var updated models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decoding updated todo: %v", err)
	}
	if updated.Title != "New title" || !updated.Completed {
		t.Errorf("updated todo = %+v", updated)
	}

	rec = doRequest(t, r, http.MethodDelete, path, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("delete status = %d", rec.Code)
	}
	rec = doRequest(t, r, http.MethodGet, path, nil)
	if rec.Code != http.StatusNotFound {
		t.Errorf("after delete, get status = %d, want 404", rec.Code)
	}
}
