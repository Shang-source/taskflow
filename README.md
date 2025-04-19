# TaskFlow

**Lightweight To‑Do List REST API built with Go and Gin**

## Overview

TaskFlow is a simple yet powerful backend service for managing to‑do tasks. It provides:

- **User Registration & Login** using JWT authentication

- **Task CRUD** (Create, Read, Update, Delete)

- **Due‑Date Reminders** via background scheduler

- **Pagination & Sorting** by due date

- **Structured Logging Middleware** for HTTP request monitoring

- **Auto‑Generated Swagger / OpenAPI Documentation**


## Prerequisites

- Go 1.20 or higher

- MySQL 5.7+ (or compatible MariaDB)

- `swag` CLI for generating Swagger docs:

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```


## Configuration

Create a `.env` file in the project root with the following variables:

```dotenv
# Database settings
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=your_password
DB_NAME=taskflow_db

# JWT settings
JWT_SECRET=YourSecretKey
JWT_EXPIRE_HOURS=24

# Server port
RUN_PORT=8080
```

## Database Initialization

On startup, the application will automatically create the required tables (users, tasks, notifications). To manually apply schema changes, run:

```sql
ALTER TABLE tasks ADD COLUMN due_date DATETIME NULL;
ALTER TABLE tasks ADD COLUMN reminded BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS notifications (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  task_id BIGINT NOT NULL,
  message VARCHAR(255) NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Build & Run

1. Install dependencies:

    ```bash
    go mod tidy
    ```

2. Generate Swagger documentation:

    ```bash
    swag init -g main.go --parseDependency
    ```

3. Start the server:

    ```bash
    go run main.go
    ```

4. Open Swagger UI in your browser:

    ```
    http://localhost:8080/swagger/index.html
    ```


## Project Structure

```
taskflow/
├─ controllers/       # API handlers (auth, tasks, notifications)
├─ middlewares/       # JWT auth & logging middleware
├─ models/            # Data models (User, Task, Notification)
├─ jobs/              # Background scheduler for reminders
├─ routes/            # HTTP route definitions
├─ config/            # Database connection & schema initialization
├─ docs/              # Generated Swagger files
├─ main.go            # Application entry point
├─ go.mod             # Go module file
└─ .env               # Environment variables
```

## API Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'
```

### Login and Obtain JWT

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'
```

### Authorize with Bearer Token

In Swagger UI, click **Authorize** and enter:

```
Bearer <your_jwt_token>
```

### Create a Task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","due_date":"2025-04-18T12:00:00Z"}'
```

## Due‑Date Reminders

- A scheduler runs every minute, selecting tasks where:

    - `due_date <= now`

    - `is_done = false`

    - `reminded = false`

- It inserts a record into the `notifications` table and sets `reminded = true`.

- Use `GET /notifications` to fetch unread reminders.


## Logging Middleware

All requests pass through `middlewares.Logger()`, which logs HTTP method, path, and processing time.

## License

This project is licensed under the MIT License. Feel free to use and modify.