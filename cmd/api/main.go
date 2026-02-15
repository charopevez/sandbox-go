// =============================================================
// Simple REST API — CRUD for Tasks
// Run: go run cmd/api/main.go
// Test: curl http://localhost:8080/tasks
//
// This is what they might ask you to build in the live coding.
// Study this pattern well — it covers:
//   - HTTP routing (net/http)
//   - JSON encode/decode
//   - Error handling
//   - Database integration
//   - Context usage
// =============================================================
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// -----------------------------------------------------------
// MODELS
// -----------------------------------------------------------

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
}

type CreateTaskRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
}

type UpdateTaskRequest struct {
	Title *string `json:"title,omitempty"` // pointer = can detect missing vs empty
	Done  *bool   `json:"done,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// -----------------------------------------------------------
// APP — holds dependencies (like a service container in PHP)
// -----------------------------------------------------------
type App struct {
	DB *pgxpool.Pool
}

// -----------------------------------------------------------
// HELPERS
// -----------------------------------------------------------

// writeJSON — helper to send JSON responses
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError — helper to send error responses
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

// extractID — get ID from URL path like /tasks/123
func extractID(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.TrimSuffix(idStr, "/")
	return strconv.Atoi(idStr)
}

// -----------------------------------------------------------
// HANDLERS
// -----------------------------------------------------------

// GET /tasks — list all tasks
func (app *App) handleListTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(r.Context(),
		"SELECT id, user_id, title, done FROM tasks ORDER BY id",
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query tasks")
		log.Printf("listTasks: %v", err)
		return
	}
	defer rows.Close()

	tasks := []Task{} // empty slice, not nil (so JSON is [] not null)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Done); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan task")
			return
		}
		tasks = append(tasks, t)
	}

	writeJSON(w, http.StatusOK, tasks)
}

// POST /tasks — create a task
func (app *App) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Validation
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.UserID == 0 {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	var task Task
	err := app.DB.QueryRow(r.Context(),
		"INSERT INTO tasks (user_id, title) VALUES ($1, $2) RETURNING id, user_id, title, done",
		req.UserID, req.Title,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Done)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		log.Printf("createTask: %v", err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

// GET /tasks/{id} — get single task
func (app *App) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/tasks/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var task Task
	err = app.DB.QueryRow(r.Context(),
		"SELECT id, user_id, title, done FROM tasks WHERE id = $1", id,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Done)

	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// PUT /tasks/{id} — update a task
func (app *App) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/tasks/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Build update dynamically (only update provided fields)
	if req.Title != nil {
		_, err = app.DB.Exec(r.Context(),
			"UPDATE tasks SET title = $1 WHERE id = $2", *req.Title, id)
	}
	if req.Done != nil {
		_, err = app.DB.Exec(r.Context(),
			"UPDATE tasks SET done = $1 WHERE id = $2", *req.Done, id)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	// Return updated task
	var task Task
	err = app.DB.QueryRow(r.Context(),
		"SELECT id, user_id, title, done FROM tasks WHERE id = $1", id,
	).Scan(&task.ID, &task.UserID, &task.Title, &task.Done)

	if err != nil {
		writeError(w, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// DELETE /tasks/{id}
func (app *App) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, "/tasks/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	tag, err := app.DB.Exec(r.Context(),
		"DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	if tag.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, fmt.Sprintf("task %d not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 — success, no body
}

// -----------------------------------------------------------
// ROUTER — simple routing without external libraries
// -----------------------------------------------------------
func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	// /tasks — collection endpoint
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.handleListTasks(w, r)
		case http.MethodPost:
			app.handleCreateTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// /tasks/{id} — single resource endpoint
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.handleGetTask(w, r)
		case http.MethodPut:
			app.handleUpdateTask(w, r)
		case http.MethodDelete:
			app.handleDeleteTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}

// -----------------------------------------------------------
// MAIN
// -----------------------------------------------------------
func main() {
	// Connect to database
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "gouser")
	pass := getEnv("DB_PASSWORD", "gopass")
	name := getEnv("DB_NAME", "sandbox")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	app := &App{DB: pool}

	// Start server
	addr := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", addr)
	fmt.Println("   GET    /tasks       — list all tasks")
	fmt.Println("   POST   /tasks       — create task")
	fmt.Println("   GET    /tasks/{id}  — get task")
	fmt.Println("   PUT    /tasks/{id}  — update task")
	fmt.Println("   DELETE /tasks/{id}  — delete task")
	fmt.Println("   GET    /health      — health check")

	log.Fatal(http.ListenAndServe(addr, app.routes()))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
