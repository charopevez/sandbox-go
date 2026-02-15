// =============================================================
// Go + PostgreSQL — CRUD Operations
// Run: go run cmd/examples/03_database.go
//
// PHP equivalent: PDO
// Go equivalent: database/sql (stdlib) or pgx (we use this)
//
// pgx is the most popular PostgreSQL driver for Go.
// It's what you'd use at Tribeat with their AWS stack.
// =============================================================
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// -----------------------------------------------------------
// MODEL — struct maps to database row
// Like a PHP entity / model class
// -----------------------------------------------------------
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

type Task struct {
	ID        int
	UserID    int
	Title     string
	Done      bool
	CreatedAt time.Time
}

// -----------------------------------------------------------
// DATABASE CONNECTION
// PHP: new PDO("pgsql:host=...", $user, $pass)
// Go:  pgxpool.New(ctx, connString)
// -----------------------------------------------------------
func connectDB() *pgxpool.Pool {
	// Build connection string from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "gouser")
	pass := getEnv("DB_PASSWORD", "gopass")
	name := getEnv("DB_NAME", "sandbox")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	// pgxpool = connection pool (like PHP's persistent connections)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect: %v\n", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping: %v\n", err)
	}

	fmt.Println("✅ Connected to database!")
	return pool
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// -----------------------------------------------------------
// SELECT MULTIPLE ROWS
// PHP: $stmt = $pdo->query("SELECT ..."); $rows = $stmt->fetchAll();
// Go:  pool.Query(ctx, "SELECT ...") → iterate with rows.Next()
// -----------------------------------------------------------
func listUsers(pool *pgxpool.Pool) ([]User, error) {
	ctx := context.Background()

	rows, err := pool.Query(ctx, "SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close() // ⚠️ ALWAYS close rows — like closing a cursor

	var users []User
	for rows.Next() {
		var u User
		// Scan maps columns to struct fields (positional, like PDO::FETCH_NUM)
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}

	// Check for errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return users, nil
}

// -----------------------------------------------------------
// SHORTCUT: pgx.CollectRows — less boilerplate
// -----------------------------------------------------------
func listUsersShort(pool *pgxpool.Pool) ([]User, error) {
	ctx := context.Background()

	rows, err := pool.Query(ctx, "SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (User, error) {
		var u User
		err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
		return u, err
	})
}

// -----------------------------------------------------------
// SELECT SINGLE ROW
// PHP: $stmt->fetch() with LIMIT 1
// Go:  pool.QueryRow() — returns exactly one row
// -----------------------------------------------------------
func getUserByID(pool *pgxpool.Pool, id int) (*User, error) {
	ctx := context.Background()

	var u User
	err := pool.QueryRow(ctx,
		"SELECT id, name, email, created_at FROM users WHERE id = $1", // $1 = parameterized
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user %d not found", id)
		}
		return nil, fmt.Errorf("query user %d: %w", id, err)
	}

	return &u, nil
}

// -----------------------------------------------------------
// INSERT — returning the new ID
// PHP: $pdo->exec("INSERT ..."); $id = $pdo->lastInsertId();
// Go:  use RETURNING clause with QueryRow
// -----------------------------------------------------------
func createUser(pool *pgxpool.Pool, name, email string) (*User, error) {
	ctx := context.Background()

	var u User
	err := pool.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at",
		name, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &u, nil
}

// -----------------------------------------------------------
// UPDATE
// PHP: $stmt = $pdo->prepare("UPDATE ..."); $stmt->execute([...]);
// -----------------------------------------------------------
func updateUserName(pool *pgxpool.Pool, id int, newName string) error {
	ctx := context.Background()

	tag, err := pool.Exec(ctx,
		"UPDATE users SET name = $1 WHERE id = $2",
		newName, id,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d not found", id)
	}

	fmt.Printf("  updated %d row(s)\n", tag.RowsAffected())
	return nil
}

// -----------------------------------------------------------
// DELETE
// -----------------------------------------------------------
func deleteUser(pool *pgxpool.Pool, id int) error {
	ctx := context.Background()

	tag, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	fmt.Printf("  deleted %d row(s)\n", tag.RowsAffected())
	return nil
}

// -----------------------------------------------------------
// TRANSACTIONS
// PHP: $pdo->beginTransaction(); ... $pdo->commit();
// Go:  pool.Begin() returns a Tx, use defer tx.Rollback()
// -----------------------------------------------------------
func createUserWithTask(pool *pgxpool.Pool, name, email, taskTitle string) error {
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// If anything fails, rollback. If we commit, rollback is a no-op.
	defer tx.Rollback(ctx)

	// Insert user within transaction
	var userID int
	err = tx.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		name, email,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	// Insert task within same transaction
	_, err = tx.Exec(ctx,
		"INSERT INTO tasks (user_id, title) VALUES ($1, $2)",
		userID, taskTitle,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Printf("  created user %d with task '%s'\n", userID, taskTitle)
	return nil
}

// -----------------------------------------------------------
// JOINS — get tasks with user info
// -----------------------------------------------------------
func getTasksWithUsers(pool *pgxpool.Pool) error {
	ctx := context.Background()

	rows, err := pool.Query(ctx, `
		SELECT t.id, t.title, t.done, u.name
		FROM tasks t
		JOIN users u ON u.id = t.user_id
		ORDER BY t.id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id       int
			title    string
			done     bool
			userName string
		)
		if err := rows.Scan(&id, &title, &done, &userName); err != nil {
			return err
		}
		status := "⬜"
		if done {
			status = "✅"
		}
		fmt.Printf("  %s [%d] %s (assigned to %s)\n", status, id, title, userName)
	}

	return rows.Err()
}

// -----------------------------------------------------------
// MAIN
// -----------------------------------------------------------
func main() {
	pool := connectDB()
	defer pool.Close()

	// List all users
	fmt.Println("\n=== LIST USERS ===")
	users, err := listUsers(pool)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		fmt.Printf("  [%d] %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// Get single user
	fmt.Println("\n=== GET USER BY ID ===")
	user, err := getUserByID(pool, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found: %s (%s)\n", user.Name, user.Email)

	// Create user
	fmt.Println("\n=== CREATE USER ===")
	newUser, err := createUser(pool, "Dave", "dave@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Created: [%d] %s\n", newUser.ID, newUser.Name)

	// Update user
	fmt.Println("\n=== UPDATE USER ===")
	if err := updateUserName(pool, newUser.ID, "David"); err != nil {
		log.Fatal(err)
	}

	// Transaction: create user + task
	fmt.Println("\n=== TRANSACTION ===")
	if err := createUserWithTask(pool, "Eve", "eve@example.com", "Join the team"); err != nil {
		log.Fatal(err)
	}

	// Join query
	fmt.Println("\n=== TASKS WITH USERS ===")
	if err := getTasksWithUsers(pool); err != nil {
		log.Fatal(err)
	}

	// Delete the test user
	fmt.Println("\n=== DELETE USER ===")
	if err := deleteUser(pool, newUser.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n✅ All database examples done!")
}
