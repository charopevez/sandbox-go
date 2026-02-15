# sandbox-go — Interview Prep Environment

Go practice environment with PostgreSQL, VS Code devcontainer, and examples
tailored for a backend engineer interview.

## Quick Start

### Prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [VS Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Setup (< 5 minutes)

```bash
# 1. Clone / copy this repo
cd sandbox-go

# 2. Open in VS Code
code .

# 3. When prompted "Reopen in Container" → click Yes
#    Or: Ctrl+Shift+P → "Dev Containers: Reopen in Container"

# 4. Wait for build (~2 min first time)
#    PostgreSQL starts automatically with seed data
```

### If NOT using devcontainers (local Go install)

```bash
# Start just the database
docker compose up db -d

# Install Go: https://go.dev/dl/
# Then:
go mod download
go run cmd/examples/01_syntax.go
```

## Project Structure

```
sandbox-go/
├── .devcontainer/         ← VS Code devcontainer config
│   └── devcontainer.json
├── cmd/
│   ├── examples/
│   │   ├── 01_syntax.go       ← Go syntax refresher (with PHP comparisons!)
│   │   ├── 02_concurrency.go  ← Goroutines, channels, worker pools
│   │   └── 03_database.go     ← PostgreSQL CRUD with pgx
│   └── api/
│       └── main.go            ← REST API server (interview-ready pattern)
├── docker-compose.yml     ← Go app + PostgreSQL
├── Dockerfile             ← Go dev container
├── init.sql               ← Database seed data
├── go.mod
└── README.md
```

## Running the Examples

From inside the devcontainer terminal:

```bash
# 1. Syntax refresher (no database needed)
go run cmd/examples/01_syntax.go

# 2. Concurrency patterns (no database needed)
go run cmd/examples/02_concurrency.go

# 3. Database operations (needs PostgreSQL running)
go run cmd/examples/03_database.go

# 4. REST API server
go run cmd/api/main.go
# Then in another terminal:
curl http://localhost:8080/tasks
curl -X POST http://localhost:8080/tasks -d '{"user_id":1,"title":"New task"}'
curl http://localhost:8080/tasks/1
curl -X PUT http://localhost:8080/tasks/1 -d '{"done":true}'
curl -X DELETE http://localhost:8080/tasks/1
```

## Study Order (6-8 hours)

### Day 1 — Today (2-3 hours)
1. ✅ Set up this environment
2. 📖 Read through `01_syntax.go` — run it, modify things, break things
3. 📖 Read through `02_concurrency.go` — this is THE Go differentiator
4. 🏋️ Try writing a small program from scratch: e.g., a word counter

### Day 2 — Tomorrow (2-3 hours)
1. 📖 Read through `03_database.go` — run it, try adding queries
2. 📖 Study `cmd/api/main.go` — this is the most likely live coding pattern
3. 🏋️ Try building the API from scratch without looking (muscle memory)
4. 🏋️ Practice: add a `/users` endpoint to the API yourself

### Day 3 — Day before interview (1-2 hours)
1. 🏋️ Build a simple API from scratch again, timed (aim for 25 min)
2. 📝 Review goroutine patterns (worker pool, fan-in/fan-out)
3. 📝 Review error handling patterns
4. 💡 Make sure you can explain: "Why Go vs PHP for this use case?"

### Day 4 — Interview day (30 min review)
1. Quick scan of key patterns
2. Make sure your laptop is ready (Go works, IDE works, docker works)
3. Open a blank project so you can start coding immediately

## Key Go Idioms to Remember

```go
// Always handle errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// defer for cleanup (like PHP's finally)
f, _ := os.Open("file.txt")
defer f.Close()

// Short if with init
if err := doThing(); err != nil {
    log.Fatal(err)
}

// Table-driven tests
tests := []struct{
    name string
    input int
    want  int
}{
    {"positive", 5, 25},
    {"zero", 0, 0},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := square(tt.input)
        if got != tt.want {
            t.Errorf("got %d, want %d", got, tt.want)
        }
    })
}
```

## Database Connection

From devcontainer or when docker-compose is running:
- **Host:** localhost (or `db` from inside devcontainer)
- **Port:** 5432
- **User:** gouser
- **Password:** gopass
- **Database:** sandbox

Connect with psql:
```bash
psql -h db -U gouser -d sandbox
# or from host:
psql -h localhost -U gouser -d sandbox
```

## Go vs PHP — Quick Mental Map

| PHP | Go |
|-----|-----|
| `$var = "hello"` | `var := "hello"` |
| `array(1,2,3)` | `[]int{1,2,3}` (slice) |
| `["k"=>"v"]` | `map[string]string{"k":"v"}` |
| `class Foo {}` | `type Foo struct {}` |
| `interface Bar {}` | `type Bar interface {}` (implicit!) |
| `try/catch` | `if err != nil {}` |
| `new PDO(...)` | `pgxpool.New(ctx, connStr)` |
| `$stmt->fetchAll()` | `rows.Next() + rows.Scan()` |
| threads/processes | goroutines (lightweight, millions ok) |
| `composer require` | `go get package` |
| `phpunit` | `go test` (built-in!) |
