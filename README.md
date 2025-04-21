# TaskFlow

TaskFlow is a simple Go-based backend for task management. It provides RESTful APIs for user registration, login (JWT-based authentication), and task retrieval, all documented with Swagger.

## Features

- User registration (`POST /register`)
- User login with JWT issuance (`POST /login`)
- Protected task list retrieval (`GET /tasks` with `Authorization` header)
- Auto-generated Swagger UI for API exploration
- Environment-based configuration
- Automatic table creation for MySQL (users & tasks)

## Prerequisites

- Go 1.18+
- MySQL 5.7+ (or compatible)
- Docker & Docker Compose (optional, for local MySQL)

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/taskflow.git
   cd taskflow
   ```

2. **Environment Variables**
   Copy `env.sample` to `.env` and update values:
   ```env
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASS=stt123456
   DB_NAME=taskflow_db

   JWT_SECRET=MySecretKey
   JWT_EXPIRE_HOURS=24

   RUN_PORT=8080
   ```

3. **Start MySQL** (optional with Docker)
   ```bash
   docker-compose up -d mysql
   ```

4. **Install dependencies**
   ```bash
   go mod tidy
   ```

## Database Migration

Tables are created automatically on startup. Alternatively, run the migration script:

```bash
go run cmd/migrate.go
```

## Running the Server

1. **Generate Swagger docs & start**
   ```bash
   swag init && go run main.go
   ```

2. **Open Swagger UI**
   Visit: http://localhost:8080/swagger/  
   Use **Try it out** to test endpoints.

## API Endpoints

### Register
- **POST** `/register`
- Body: `{ "username": "user1", "password": "pass" }`
- Response: `{ "message": "Registration successful", "user": { "id": 1, "username": "user1" } }`

### Login
- **POST** `/login`
- Body: `{ "username": "user1", "password": "pass" }`
- Response: `{ "message": "Login successful", "token": "<JWT>" }`

### Get Tasks
- **GET** `/tasks`
- Header: `Authorization: Bearer <JWT>`
- Response: `[{ "id": 1, "title": "Task 1", ... }]`

## Project Structure

```
taskflow/
├── cmd/               # Migration script
│   └── migrate.go
├── config/            # Environment & DB connection
│   └── database.go
├── controllers/       # HTTP handlers
│   ├── auth.go
│   └── task.go
├── middlewares/       # Gin middleware (logging, auth)
├── models/            # Data models & Swagger annotations
├── routes/            # URL routing setup
├── docs/              # Generated Swagger files
├── .env               # Environment variables
├── main.go            # Application entry point
├── go.mod             # Go module definition
└── go.sum             # Dependency checksums
```

