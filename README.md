# Todo App (Go)

A todo list with a Go + Gin API and a small web page. Tasks are
stored in a SQLite file — no database server to install.

## Run it

You need Go 1.25 or newer.

```sh
go run .
```

Then open <http://localhost:8080> in your browser. The page talks
to the API on the same server, so it works on any host and port.

Change the port or the database file with environment variables:

```sh
PORT=9090 TODO_DB_PATH=my-tasks.db go run .
```

See `.env.example` for the settings. The database file
(`todo.db`) is created automatically on the first run, and the
`todos` table is created automatically too — see
`database/schema.sql` for the shape of the table.

## API

| Method | Path | What it does |
| --- | --- | --- |
| GET | `/api/health` | Says the API is running |
| POST | `/api/todos` | Create a todo (title is required) |
| GET | `/api/todos` | List all todos |
| GET | `/api/todos/:id` | Get one todo, or 404 |
| PUT | `/api/todos/:id` | Replace a todo, or 404 |
| DELETE | `/api/todos/:id` | Delete a todo, or 404 |

A todo looks like this:

```json
{
  "ID": 1,
  "title": "Buy milk",
  "description": "From the shop on the corner",
  "completed": false,
  "priority": "high",
  "due_date": "2026-11-20T00:00:00Z",
  "category": "shopping"
}
```

Priority must be `low`, `medium` or `high`. Category must be
`work`, `personal`, `shopping`, `health` or `other`.

A ready-made Postman collection is in
`postman/todo-app-go.postman_collection.json`.

## Tests

```sh
go test ./...
go vet ./...
```

The tests run the real handlers against an in-memory SQLite
database: create, list, get, update, delete, and validation
(empty title, bad priority, missing todo).

## Files

- `main.go` — server setup, routes, port and graceful shutdown
- `handlers/todo_handler.go` — all five todo endpoints (this file
  used to be empty; the logic used to live in `main.go`)
- `models/todo.go` — the Todo type
- `database/db.go` — SQLite connection and automatic migration
- `database/schema.sql` — the table, for reference
- `public/` — the web page (`index.html`, `app.js`, `style.css`)
- `postman/` — a Postman collection for the API

## Note

The old `todo-app.exe` (a 32 MB Windows binary) was removed from
the repository — binaries do not belong in git. Build your own
with `go build .`, or download one from a GitHub Release.
