# Task Runner API

A simple in-memory HTTP service for running long tasks.

---

## Features

- Create tasks by type (e.g. "default")
- Check task status, result, and duration
- Delete tasks (except if running)
- One task runs at a time for each task type
- No database, queues, or external services

---

## How to Run

```bash
go run ./cmd/task-runner
```

Server will start on: `http://localhost:8080`

---

## API

### Create Task

```
POST /tasks?type=default
```

**Response:**

```json
{
  "id": "abc123...",
  "status": "pending",
  "created_at": "2025-06-19T12:00:00Z"
}
```

---

### Get Task by ID

```
GET /tasks/{id}
```

**Response:**

```json
{
  "id": "abc123...",
  "status": "running",
  "created_at": "2025-06-19T12:00:00Z",
  "duration": "00:00:12",
  "result": ""
}
```

---

### Delete Task

```
DELETE /tasks/{id}
```

**Responses:**

- `204 No Content` — success
- `404 Not Found` — task not found
- `409 Conflict` — task is still running

---

## Project Structure

```
cmd/                  # Entry point
internal/bootstrap/   # Task type registration
internal/domain/      # Task manager and task logic
internal/transport/   # HTTP API
```

---

## Add New Task Types

1. Implement the `ExecutableTask` interface
2. Add a factory that creates the task
3. Register it in `RegisterTaskFactories(...)`

---

## Requirements

- Go 1.20+
- No external services
